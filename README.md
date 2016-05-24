# gofilemon
Tool to monitor files for specific regular expressions, and exit success or failure depending on what is found.

Usage:
  gofilemon [OPTIONS]

Application Options:
  -s, --succeed= Monitor the file for the supplied regex, and exit with success if found.  Use the following format: <file to
                 monitor>:<regex to look for>
  -f, --fail=    Monitor the file for the supplied regex, and exit with failure if found.  Use the following format: <file to
                 monitor>:<regex to look for>

Help Options:
  -h, --help     Show this help message

This command will monitor files for specific regular expressions, and, depending on the options provided, will exit with either a successful return code, or a failure return code when those regular expressions are found.  The file does not need to exist prior to running this
command.

