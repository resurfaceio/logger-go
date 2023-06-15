package logger

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var bodies = map[string][]byte{
	"plain":   {10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 46, 45, 39, 39, 39, 39, 39, 39, 45, 46, 10, 32, 32, 32, 32, 32, 32, 32, 32, 46, 39, 32, 95, 32, 32, 32, 32, 32, 32, 95, 32, 39, 46, 10, 32, 32, 32, 32, 32, 32, 32, 47, 32, 32, 32, 79, 32, 32, 32, 32, 32, 32, 79, 32, 32, 32, 92, 10, 32, 32, 32, 32, 32, 32, 58, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 58, 10, 32, 32, 32, 32, 32, 32, 124, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 124, 10, 32, 32, 32, 32, 32, 32, 58, 32, 32, 32, 32, 32, 32, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 58, 10, 32, 32, 32, 32, 32, 32, 32, 92, 32, 32, 46, 45, 34, 96, 32, 32, 96, 34, 45, 46, 32, 32, 47, 10, 32, 32, 32, 32, 32, 32, 32, 32, 39, 46, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 46, 39, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 39, 45, 46, 46, 46, 46, 46, 46, 45, 39, 10, 32, 32, 32, 32, 32, 89, 79, 85, 32, 83, 72, 79, 85, 76, 68, 78, 39, 84, 32, 66, 69, 32, 72, 69, 82, 69, 10},
	"json":    {123, 10, 32, 32, 34, 115, 108, 105, 100, 101, 115, 104, 111, 119, 34, 58, 32, 123, 10, 32, 32, 32, 32, 34, 97, 117, 116, 104, 111, 114, 34, 58, 32, 34, 89, 111, 117, 114, 115, 32, 84, 114, 117, 108, 121, 34, 44, 32, 10, 32, 32, 32, 32, 34, 100, 97, 116, 101, 34, 58, 32, 34, 100, 97, 116, 101, 32, 111, 102, 32, 112, 117, 98, 108, 105, 99, 97, 116, 105, 111, 110, 34, 44, 32, 10, 32, 32, 32, 32, 34, 115, 108, 105, 100, 101, 115, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 105, 116, 108, 101, 34, 58, 32, 34, 87, 97, 107, 101, 32, 117, 112, 32, 116, 111, 32, 87, 111, 110, 100, 101, 114, 87, 105, 100, 103, 101, 116, 115, 33, 34, 44, 32, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 121, 112, 101, 34, 58, 32, 34, 97, 108, 108, 34, 10, 32, 32, 32, 32, 32, 32, 125, 44, 32, 10, 32, 32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 105, 116, 101, 109, 115, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 87, 104, 121, 32, 60, 101, 109, 62, 87, 111, 110, 100, 101, 114, 87, 105, 100, 103, 101, 116, 115, 60, 47, 101, 109, 62, 32, 97, 114, 101, 32, 103, 114, 101, 97, 116, 34, 44, 32, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 87, 104, 111, 32, 60, 101, 109, 62, 98, 117, 121, 115, 60, 47, 101, 109, 62, 32, 87, 111, 110, 100, 101, 114, 87, 105, 100, 103, 101, 116, 115, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 93, 44, 32, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 105, 116, 108, 101, 34, 58, 32, 34, 79, 118, 101, 114, 118, 105, 101, 119, 34, 44, 32, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 121, 112, 101, 34, 58, 32, 34, 97, 108, 108, 34, 10, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 93, 44, 32, 10, 32, 32, 32, 32, 34, 116, 105, 116, 108, 101, 34, 58, 32, 34, 83, 97, 109, 112, 108, 101, 32, 83, 108, 105, 100, 101, 32, 83, 104, 111, 119, 34, 10, 32, 32, 125, 10, 125, 10},
	"gzip":    {31, 139, 8, 0, 130, 139, 110, 100, 2, 255, 61, 143, 65, 14, 131, 32, 20, 68, 247, 158, 130, 176, 52, 138, 232, 198, 214, 157, 105, 76, 123, 128, 246, 0, 20, 127, 149, 212, 2, 5, 108, 210, 26, 239, 94, 192, 196, 229, 204, 188, 204, 100, 150, 4, 33, 60, 252, 132, 214, 208, 227, 6, 57, 51, 67, 134, 130, 55, 2, 235, 193, 88, 239, 45, 94, 122, 163, 229, 28, 180, 243, 26, 167, 69, 138, 35, 180, 187, 121, 39, 185, 234, 133, 28, 66, 28, 218, 50, 212, 195, 99, 98, 206, 151, 221, 205, 14, 159, 148, 148, 192, 157, 80, 50, 112, 79, 0, 157, 179, 73, 124, 96, 7, 46, 202, 198, 133, 73, 113, 54, 141, 94, 52, 7, 74, 233, 30, 223, 44, 152, 188, 29, 64, 70, 72, 127, 221, 168, 100, 110, 224, 61, 131, 117, 182, 168, 72, 117, 36, 20, 123, 116, 221, 46, 188, 192, 3, 225, 21, 62, 119, 215, 173, 4, 43, 35, 6, 17, 231, 203, 186, 34, 101, 77, 40, 41, 113, 178, 38, 127, 190, 200, 245, 32, 8, 1, 0, 0},
	"deflate": {120, 156, 61, 143, 65, 14, 130, 48, 20, 68, 247, 156, 162, 233, 146, 64, 41, 108, 80, 119, 196, 16, 61, 128, 30, 160, 150, 47, 52, 214, 182, 182, 197, 68, 9, 119, 183, 5, 195, 114, 102, 94, 102, 50, 83, 130, 16, 238, 224, 46, 153, 135, 14, 31, 144, 183, 35, 100, 40, 154, 3, 176, 14, 172, 11, 222, 20, 100, 48, 26, 206, 193, 248, 160, 113, 90, 164, 120, 129, 54, 55, 111, 21, 215, 157, 80, 125, 140, 251, 175, 48, 25, 250, 151, 102, 232, 102, 55, 248, 168, 149, 2, 238, 133, 86, 145, 123, 0, 152, 156, 73, 241, 134, 13, 56, 107, 183, 44, 72, 205, 153, 28, 130, 56, 236, 40, 165, 91, 124, 117, 96, 243, 166, 7, 181, 64, 230, 227, 7, 173, 114, 11, 175, 17, 156, 119, 69, 69, 170, 61, 161, 56, 160, 243, 122, 225, 9, 1, 136, 175, 240, 169, 189, 172, 37, 88, 91, 209, 139, 101, 190, 172, 43, 82, 214, 132, 146, 18, 39, 115, 242, 3, 1, 84, 71, 126},
	"br":      {27, 6, 1, 0, 156, 7, 182, 227, 165, 196, 166, 120, 3, 133, 227, 113, 62, 253, 57, 242, 135, 63, 114, 13, 52, 144, 157, 71, 69, 130, 136, 244, 133, 66, 232, 88, 238, 182, 128, 2, 12, 136, 175, 219, 178, 246, 123, 105, 251, 162, 215, 32, 208, 84, 243, 74, 231, 22, 8, 36, 176, 48, 153, 114, 8, 6, 45, 204, 231, 79, 217, 54, 31, 67, 139, 115, 44, 204, 204, 167, 239, 48, 107, 208, 87, 140, 56, 1, 145, 4, 204, 168, 252, 194, 24, 52, 213, 95, 187, 88, 81, 86, 207, 40, 91, 182, 57, 108, 141, 136, 48, 64, 73, 121, 186, 185, 34, 61, 158, 247, 85, 114, 4, 192, 104, 69, 129, 82, 138, 102, 178, 151, 27, 239, 2, 197, 7, 150, 247, 104, 230, 137, 87, 224, 85, 26, 97, 66, 161, 192, 97, 16, 47, 7, 212, 34, 237, 204, 82, 16, 237, 59, 69, 237, 139, 208, 172, 3},
	// "compress": {},
	// "zstd":     {},
	"jpg":        {255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 1, 0, 72, 0, 72, 0, 0, 255, 219, 0, 67, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 219, 0, 67, 1, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 192, 0, 17, 8, 0, 234, 1, 57, 3, 1, 34, 0, 2, 17, 1, 3, 17, 1, 255, 196, 0, 23, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 255, 196, 0, 36, 16, 1, 1, 1, 0, 2, 1, 4, 3, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 17, 33, 49, 65, 2, 18, 81, 113, 97, 129, 145, 177, 161, 193, 240, 255, 196, 0, 21, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 255, 196, 0, 22, 17, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 17, 1, 255, 218, 0, 12, 3, 1, 0, 2, 17, 3, 17, 0, 63, 0, 202, 226, 47, 32, 138, 130, 13, 96, 67, 0, 160, 179, 129, 83, 48, 231, 110, 181, 234, 99, 65, 169, 62, 23, 51, 18, 86, 168, 49, 101, 212, 93, 60, 242, 8, 188, 126, 86, 207, 49, 157, 6, 184, 169, 134, 212, 222, 1, 120, 54, 102, 50, 128, 214, 154, 138, 13, 65, 153, 113, 173, 223, 160, 23, 164, 242, 1, 111, 240, 135, 183, 124, 167, 64, 214, 166, 252, 179, 184, 104, 47, 167, 139, 86, 222, 117, 22, 192, 39, 170, 93, 212, 183, 83, 59, 1, 102, 223, 148, 179, 3, 144, 67, 23, 252, 94, 191, 98, 50, 191, 164, 0, 197, 64, 5, 13, 3, 166, 165, 97, 168, 10, 147, 189, 166, 254, 63, 105, 224, 26, 182, 124, 48, 210, 96, 162, 205, 92, 243, 75, 160, 113, 63, 53, 62, 145, 96, 45, 245, 112, 156, 89, 223, 37, 149, 51, 230, 129, 98, 47, 210, 217, 152, 34, 2, 201, 216, 39, 30, 67, 219, 91, 207, 154, 170, 230, 220, 244, 223, 165, 4, 48, 200, 0, 188, 37, 155, 210, 40, 57, 226, 227, 86, 107, 50, 229, 21, 86, 95, 148, 207, 220, 61, 191, 253, 80, 90, 182, 108, 103, 153, 197, 60, 118, 7, 248, 150, 114, 186, 183, 144, 72, 94, 76, 249, 242, 116, 8, 7, 240, 68, 5, 81, 5, 0, 53, 12, 5, 149, 81, 172, 168, 37, 69, 168, 14, 155, 49, 158, 211, 236, 185, 224, 85, 227, 229, 113, 133, 220, 128, 182, 255, 0, 140, 226, 198, 164, 128, 204, 45, 183, 134, 189, 177, 120, 138, 140, 201, 126, 191, 214, 186, 64, 21, 0, 16, 23, 1, 5, 196, 0, 0, 18, 205, 80, 25, 116, 237, 150, 54, 138, 221, 155, 231, 166, 22, 94, 215, 113, 4, 141, 102, 33, 109, 191, 64, 189, 195, 221, 39, 132, 151, 248, 152, 13, 110, 248, 103, 9, 23, 2, 50, 2, 162, 255, 0, 234, 47, 194, 1, 141, 125, 36, 106, 116, 138, 205, 107, 111, 76, 210, 2, 130, 2, 136, 212, 244, 240, 12, 141, 101, 233, 114, 65, 9, 20, 209, 68, 0, 5, 196, 80, 16, 1, 99, 62, 171, 103, 75, 171, 121, 7, 61, 191, 43, 125, 95, 16, 197, 217, 0, 23, 180, 0, 92, 64, 11, 52, 1, 156, 176, 146, 249, 110, 37, 188, 162, 166, 39, 218, 179, 157, 131, 115, 44, 99, 149, 156, 30, 0, 237, 49, 174, 47, 225, 115, 236, 25, 56, 66, 136, 187, 169, 210, 201, 228, 232, 85, 156, 150, 95, 11, 194, 131, 157, 224, 91, 205, 80, 67, 182, 147, 52, 9, 249, 116, 102, 73, 62, 205, 84, 93, 68, 80, 1, 65, 12, 9, 116, 11, 100, 79, 119, 61, 22, 121, 132, 244, 224, 52, 202, 219, 35, 50, 130, 155, 157, 139, 103, 0, 172, 251, 99, 60, 195, 104, 55, 178, 51, 238, 100, 5, 219, 224, 139, 12, 1, 64, 2, 205, 85, 6, 51, 143, 202, 250, 103, 122, 188, 69, 212, 83, 35, 57, 60, 52, 116, 35, 22, 30, 234, 182, 166, 138, 29, 21, 52, 26, 210, 227, 26, 42, 55, 196, 53, 4, 84, 36, 223, 42, 178, 2, 200, 168, 138, 128, 0, 42, 0, 210, 36, 190, 63, 141, 2, 39, 230, 40, 11, 172, 250, 173, 58, 231, 250, 160, 230, 177, 108, 202, 0, 150, 213, 0, 134, 18, 40, 38, 40, 160, 130, 179, 238, 248, 6, 142, 35, 27, 64, 95, 114, 91, 104, 2, 26, 47, 8, 184, 183, 213, 248, 79, 119, 224, 216, 1, 109, 169, 202, 212, 81, 170, 207, 219, 76, 212, 12, 9, 170, 160, 77, 22, 68, 22, 79, 53, 85, 21, 16, 0, 0, 1, 80, 2, 242, 75, 226, 135, 96, 210, 43, 19, 157, 6, 171, 59, 103, 73, 116, 128, 141, 38, 40, 0, 184, 2, 225, 196, 237, 155, 65, 174, 35, 62, 239, 142, 16, 0, 0, 0, 0, 16, 5, 178, 247, 136, 182, 208, 72, 11, 38, 162, 146, 252, 181, 250, 191, 194, 122, 126, 107, 89, 1, 200, 5, 6, 184, 172, 128, 209, 169, 41, 121, 168, 53, 170, 206, 42, 160, 0, 0, 0, 76, 14, 193, 164, 103, 106, 93, 189, 130, 219, 226, 127, 82, 0, 45, 69, 0, 82, 66, 217, 1, 89, 190, 175, 132, 182, 208, 0, 0, 0, 0, 64, 81, 0, 0, 0, 1, 70, 230, 248, 102, 77, 116, 252, 64, 63, 234, 165, 217, 56, 240, 199, 191, 240, 34, 6, 124, 136, 171, 193, 147, 224, 137, 160, 37, 226, 139, 216, 26, 169, 159, 234, 170, 40, 64, 0, 1, 21, 0, 17, 164, 4, 85, 48, 5, 251, 77, 145, 157, 208, 91, 126, 16, 0, 0, 1, 20, 4, 0, 84, 0, 0, 0, 0, 17, 86, 73, 124, 130, 207, 248, 222, 201, 53, 140, 190, 57, 46, 241, 192, 173, 234, 100, 248, 79, 87, 76, 109, 249, 160, 186, 172, 152, 130, 166, 44, 1, 50, 128, 168, 52, 203, 80, 5, 64, 21, 23, 194, 0, 0, 2, 128, 66, 223, 132, 183, 80, 0, 0, 4, 5, 64, 0, 0, 0, 0, 16, 20, 64, 5, 64, 5, 212, 80, 107, 83, 220, 99, 40, 173, 219, 62, 124, 51, 192, 184, 162, 227, 54, 53, 122, 78, 16, 103, 150, 180, 214, 84, 106, 216, 128, 34, 128, 10, 34, 128, 168, 2, 129, 184, 11, 211, 22, 233, 216, 8, 160, 2, 0, 0, 0, 0, 0, 0, 138, 88, 8, 11, 208, 32, 0, 0, 2, 226, 0, 208, 138, 138, 131, 75, 128, 205, 66, 172, 184, 161, 153, 54, 159, 146, 221, 68, 23, 116, 153, 229, 149, 84, 91, 137, 167, 32, 171, 192, 139, 151, 225, 17, 68, 85, 13, 101, 80, 5, 16, 0, 0, 0, 0, 0, 4, 5, 69, 92, 128, 202, 233, 96, 2, 160, 138, 148, 95, 8, 168, 187, 4, 128, 47, 102, 26, 128, 190, 34, 235, 42, 10, 187, 249, 100, 69, 21, 58, 77, 81, 184, 151, 240, 72, 179, 142, 144, 97, 98, 225, 115, 194, 161, 134, 43, 54, 34, 154, 186, 139, 1, 26, 47, 72, 10, 203, 76, 213, 64, 0, 0, 0, 0, 21, 13, 21, 113, 48, 64, 93, 65, 160, 34, 161, 168, 42, 18, 128, 205, 20, 84, 65, 64, 48, 198, 132, 171, 25, 26, 100, 64, 15, 234, 138, 139, 14, 47, 218, 9, 166, 211, 17, 69, 86, 90, 209, 85, 42, 74, 32, 161, 19, 160, 106, 50, 0, 176, 168, 170, 136, 162, 0, 8, 2, 162, 130, 42, 2, 168, 162, 8, 52, 148, 17, 110, 32, 168, 0, 0, 2, 141, 70, 87, 246, 130, 171, 31, 182, 193, 45, 101, 106, 0, 2, 160, 0, 8, 170, 12, 168, 0, 30, 74, 138, 169, 223, 217, 11, 216, 10, 139, 58, 81, 60, 172, 64, 69, 64, 4, 88, 0, 184, 98, 130, 179, 66, 246, 1, 166, 162, 130, 233, 225, 60, 30, 17, 1, 69, 4, 80, 4, 10, 2, 40, 2, 237, 34, 124, 128, 42, 0, 0, 63, 255, 217},
	"misleading": {31, 139, 2, 0, 156, 7, 182, 211, 180, 36, 77, 218, 248, 20, 157, 222, 184, 106, 113, 78, 101, 243, 24, 125, 32, 137, 164, 50, 78, 203, 75, 123, 108, 16, 145, 174, 190, 201, 162, 160, 147, 32, 10, 101, 83, 231, 206, 31, 100, 96, 216, 165, 21, 24, 69, 21, 34, 107, 209, 233, 192, 39, 56, 157, 186, 199, 160, 207, 206, 22, 38, 176, 68, 53, 184, 197, 247, 127, 180, 242, 125, 109, 113, 173, 66, 28, 36, 219, 91, 255, 127, 132, 19, 144, 173, 19, 14, 75, 239, 106, 160, 138, 3, 75, 201, 174, 133, 218, 112, 66, 65, 75, 239, 234, 18, 234, 4, 150, 193, 152, 204, 225, 174, 247, 236, 128, 144, 68, 136, 78, 172, 212, 128, 150, 6, 168, 236, 68, 144, 31, 137, 11, 65, 106, 75, 252, 182, 232, 73, 128, 43, 176, 212, 58, 2, 67, 101, 22, 28, 42, 83, 46, 147, 239, 65, 49, 234, 89, 178, 236, 48, 23, 21, 238, 117, 118, 110, 28, 226, 202, 55, 234, 116, 207, 77, 60, 80, 40, 31, 251, 22, 132, 147, 21, 191, 117, 132, 129, 70, 18, 21, 226, 209, 166, 159, 66, 130, 71, 26, 70, 14, 240, 20, 88, 199, 112, 44, 252, 168, 17, 232, 74, 154, 154, 71, 30, 159, 235, 138, 152, 248, 236, 62, 8, 21, 237, 27, 182, 128, 140, 65, 100, 227, 152, 70, 7, 96, 191, 87, 28, 120, 188, 118, 95, 24, 110, 199, 125, 96, 87, 123, 143, 250, 26, 70, 232, 10, 162, 6, 188, 81, 228, 47, 123, 104, 64, 194, 182, 53, 28, 35, 11, 205, 150, 249, 129, 9, 55, 230, 32, 48, 120, 65, 15, 246, 22, 142, 61, 166, 134, 77, 65, 125},
}

func mockGetCustomRequest(reqURL string, customMethod string, customHeaders map[string][]string, customBody io.ReadCloser) http.Request {
	request := http.Request{
		Method:     "GET",
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Host:       "testing.dev",
	}

	if customMethod != "" {
		request.Method = customMethod
	}

	parsedURL, err := url.Parse(reqURL)
	if err != nil {
		log.Panicln(err)
	}
	request.URL = parsedURL
	request.RequestURI = reqURL

	request.Header = http.Header{
		"User-Agent":      {"python-requests/2.29.0"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
	}

	for k, v := range customHeaders {
		joined := strings.Join(v, ",")
		request.Header.Add(k, joined)
	}

	request.Body = customBody

	return request
}

func mockGetCustomResponse(request *http.Request, customHeaders map[string][]string, customBodyBytes []byte) http.Response {
	response := http.Response{
		Request:    request,
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	response.Header = http.Header{
		"Server":                           {"gunicorn/19.9.0"},
		"Date":                             {"Wed, 24 May 2023 14:49:25 GMT"},
		"Connection":                       {"keep-alive"},
		"Access-Control-Allow-Origin":      {"*"},
		"Access-Control-Allow-Credentials": {"true"},
	}

	for k, v := range customHeaders {
		joined := strings.Join(v, ",")
		response.Header.Add(k, joined)
	}

	contentLen, err := strconv.ParseInt(strings.Join(response.Header["Content-Length"], ","), 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	response.Body = io.NopCloser(io.LimitReader(bufio.NewReader(bytes.NewReader(customBodyBytes)), contentLen))
	response.ContentLength = int64(contentLen)

	return response

}

func MockGetNoBodyRequest() http.Request {
	return mockGetCustomRequest("/json", "", nil, http.NoBody)
}

func MockGetRequestWithBody(body []byte, mimeType string) http.Request {
	headers := map[string][]string{
		"Content-Length": {strconv.Itoa(len(body))},
		"Content-Type":   {mimeType},
	}
	return mockGetCustomRequest("/post", "POST", headers, io.NopCloser(bytes.NewReader(body)))
}

func MockGetNoBodyResponse(request *http.Request) http.Response {
	return mockGetCustomResponse(request, nil, []byte{})
}

func MockGetResponseWithBody(request *http.Request, body []byte, mimeType string) http.Response {
	headers := map[string][]string{
		"Content-Length": {strconv.Itoa(len(body))},
		"Content-Type":   {mimeType},
	}
	return mockGetCustomResponse(request, headers, body)
}

func MockGetPlainTextResponse(request *http.Request) http.Response {
	return MockGetResponseWithBody(request, bodies["plain"], "text/plain")
}

func MockGetJSONResponse(request *http.Request) http.Response {
	return MockGetResponseWithBody(request, bodies["json"], "application/json")
}

func MockGetJPEGResponse(request *http.Request) http.Response {
	return MockGetResponseWithBody(request, bodies["jpg"], "image/jpeg")
}

func MockGetJSONEncodedResponse(request *http.Request, encoding string) http.Response {
	body := bodies[encoding]
	headers := map[string][]string{
		"Content-Length":   {strconv.Itoa(len(body))},
		"Content-Type":     {"application/json"},
		"Content-Encoding": {encoding},
	}
	return mockGetCustomResponse(request, headers, body)
}

func MockGetGzipResponse(request *http.Request) http.Response {
	return MockGetJSONEncodedResponse(request, "gzip")
}

func MockGetDeflateResponse(request *http.Request) http.Response {
	return MockGetJSONEncodedResponse(request, "deflate")
}

func MockGetBrotliResponse(request *http.Request) http.Response {
	return MockGetJSONEncodedResponse(request, "br")
}
