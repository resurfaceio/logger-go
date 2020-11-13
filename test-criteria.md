Contributors: Tim Lenahan

Last Updated: 11/04/2020
# Unit Test Criteria
 
## Primary Class Tests

### Base Logger
- Create basic basic instance (no constructor arguments) and check that all attributes have correct value

- Multiple instances
    - Create multiple instances with logger properly ID and valid urls attached
        - check attributes are correct
    - Disable all instances with UsageLoggers
        - check that all instances are disabled
    - Enable all instances with UsageLoggoers
        - check that all instances are enabled

- Check that the Base Logger has a valid host (for client I a assume this would check for a valid host for the api that it is hooked into?)

- Check if metadataId initialzed and properly formatted

- Check if version properly attached and formatted 

- Check if enabling functions correctly
    - Note: aside from enable method there maybe class contructors that enable the logger by default

- Check that the baselogger remains disabled when given an invalid url (api endpoint) when initialized

- Check that the baselogger remains disabled when given an null url (api endpoint) when initialized

- Test a valid submission to https url (api endpoint)
    - logger should have 0 failed submissions
    - logger should have 1 successful submission

- Test a valid submission to https url (api endpoint)
    - logger should have 0 failed submissions
    - logger should have 1 successful submission

- Test a valid submission to url (api endpoint) with skipped compression
    - logger should have 0 failed submissions
    - logger should have 1 successful submission

- Test submission to an invalid url (api endpoint)
    - logger should have 1 failed submissions
    - logger should have 0 successful submission

- Test Queue submission with empty queue
    - logger should have 0 failed submissions
    - logger should have 2 successful submission

- Test logger with skipped options
    - initialize with proper id and url
    - enable/disable options
        - check that options are properly set


# (Potential) Helper Class 

- Attributes
    - valid and invalid request bodies
    - valid and invalid urls

- Methods
    - valid and invalid requests
    - functioning and non-functioning apis