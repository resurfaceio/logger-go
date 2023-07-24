package logger

import (
	"bytes"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
	"github.com/liamg/magic"
)

const brotliFallbackEnabled bool = false

func checkMagicBytes(reader *bytes.Reader) (string, error) {
	var encoding string

	payload, err := io.ReadAll(reader)
	if err != nil {
		return "", nil
	}
	reader.Seek(0, io.SeekStart)

	// Another way:
	// encoding := http.DetectContentType(b)

	filetype, err := magic.Lookup(payload)
	if err != nil {
		// "Unknown" file, i.e. identity, brotli, something else not on this list:
		// https://github.com/liamg/magic/blob/master/types.go
		if brotliFallbackEnabled {
			return "br", err
		}
		return "", err
	}
	switch extension := filetype.Extension; extension {
	case "gz":
		encoding = "gzip"
	case "zlib":
		encoding = "deflate"
	case "":
		if !strings.Contains(filetype.Description, "UTF-8") {
			// Any file in the list (see comment above) that doesn't have
			// an extension, and is not UTF-8 with BOM, i.e. not to be read
			return filetype.Description, magic.ErrUnknown
		}
	default:
		// Any other file on the list with a known extension, i.e. not to be read
		return extension, magic.ErrUnknown
	}

	return encoding, nil
}

func checkMagic(reader *io.Reader, magicReader *bytes.Reader) (string, error) {
	// According to https://github.com/liamg/magic/blob/master/magic.go#L18 the first 1024 bytes should be provided (i.e. magicLimit = 1024)
	// However, this is not actually required anywhere in the code. Logically, magicLimit > offset + magic bytes
	// Offset is 0 for both gzip and zlib, and the number of magic bytes is 4 for both as well. Then, magicLimit = 4 if just checking for gzip/zlib
	// Largest possible offset in magic lib (excluding .iso filetype) is 257 bytes, and largest possible sum (offset + magic bytes) is 265. Thus, magicLimit = 265
	const magicLimit = 265
	magicBytes := make([]byte, magicLimit)
	n, err := (*reader).Read(magicBytes)
	if n != 0 {
		if magicReader == nil {
			magicReader = bytes.NewReader(magicBytes[:n])
		} else {
			magicReader.Reset(magicBytes[:n])
		}
	}

	if err == nil {
		if _, ok := (*reader).(*bytes.Reader); ok {
			(*reader).(*bytes.Reader).Seek(0, io.SeekStart)
		} else {
			*reader = io.MultiReader(magicReader, *reader)
		}
	} else if err == io.EOF {
		*reader, err = magicReader, nil
	} else {
		return "", err
	}

	// Get encoding from magic bytes
	magicEncoding, err := checkMagicBytes(magicReader)
	return magicEncoding, err
}

func newWrap(reader *io.Reader, magicReader *bytes.Reader, encoding string) (io.Reader, error) {
	var newReader io.Reader
	var err error

	switch encoding {
	case "gzip", "x-gzip":
		newReader, err = gzip.NewReader(magicReader)
		magicReader.Seek(0, io.SeekStart)
		if err == nil {
			gzReader, gzErr := gzip.NewReader(*reader)
			gzReader.Multistream(false)
			newReader, err = gzReader, gzErr
		}
	case "deflate", "zlib", "deflated":
		_, err = zlib.NewReader(magicReader)
		magicReader.Seek(0, io.SeekStart)
		if err == nil {
			newReader, err = zlib.NewReader(*reader)
		}
	case "br":
		// Read
		allBytes, err := io.ReadAll(*reader)
		magicReader.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		if _, ok := (*reader).(*bytes.Reader); !ok {
			*reader = bytes.NewReader(allBytes)
		} else {
			(*reader).(*bytes.Reader).Seek(0, io.SeekStart)
		}
		// Decode
		decodedBytes, err := io.ReadAll(brotli.NewReader(*reader))
		if err != nil {
			(*reader).(*bytes.Reader).Seek(0, io.SeekStart)
			return *reader, err
		}
		// Return
		(*reader).(*bytes.Reader).Reset(decodedBytes)
		newReader = *reader
	case "", "identity":
		newReader = *reader
	default:
		return nil, io.ErrNoProgress
	}
	return newReader, err
}

func wrapReader(reader *io.Reader, magicReader *bytes.Reader, encoding string, magicEncoding string) (io.Reader, error, bool) {
	var accurateMagic = true
	// First wrapping attempt assumes Content-Encoding header value is correct
	wrappedReader, err := newWrap(reader, magicReader, encoding)
	// First wrap failed (header error -> wrong type, or decoding error -> corrupted file)
	if err != nil {
		// Brotli reader check requires reading the entire thing at least once. Then, this is the same as io.ReadAll(reader) failing
		if encoding == "br" && wrappedReader == nil {
			return nil, err, accurateMagic
		}
		// Second wrapping attempt, this time with magic encoding
		if magicEncoding != encoding {
			wrappedReader, err = newWrap(reader, magicReader, magicEncoding)
			// Second wrap failed (corrupted file/magic bytes were wrong)
			if err != nil {
				accurateMagic = false
				// Fallback to brotli if magic bytes were wrong
				if encoding != "br" {
					wrappedReader, err = newWrap(reader, magicReader, "br")
					if wrappedReader == nil {
						return nil, err, accurateMagic
					}
				}
				wrappedReader, err = newWrap(reader, magicReader, "")
			}
		} else {
			accurateMagic = false
			// Fall back to brotli since magic bytes do correspond to encoding header but something went wrong
			if encoding != "br" {
				wrappedReader, err = newWrap(reader, magicReader, "br")
				if wrappedReader == nil {
					return nil, err, accurateMagic
				}
			}
			// Brotli reader check requires reading the entire thing, and then decoding it.
			// This is the decoding step failing with wrappedReader being the non-decoded bytes read.
			// (when would this happen? e.g. header says br, but it is actually identity -> err != nil, reader stays encoded)
			return wrappedReader, nil, accurateMagic
		}
	} else if magicEncoding != encoding {
		// First wrap didn't fail (all good, or misconfigured Content-Encoding header as "identity" or "br")
		wrappedReader, err = newWrap(reader, magicReader, magicEncoding)
		// Second wrap failed (corrupted file/magic bytes were wrong)
		if err != nil {
			// Fallback to brotli if magic bytes were wrong
			accurateMagic = false
			if encoding != "br" {
				wrappedReader, err = newWrap(reader, magicReader, "br")
				if wrappedReader == nil {
					return nil, err, accurateMagic
				}
			}
			wrappedReader, err = newWrap(reader, magicReader, "")
		}
	}

	return wrappedReader, err, accurateMagic
}
