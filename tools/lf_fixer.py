import os
import csv
import argparse
import datetime
import typing
import fnmatch

input_dir_arg = "--input"
output_dir_arg = "--output"
files_selector_arg = "--selector"

parser=argparse.ArgumentParser()
parser.add_argument(input_dir_arg, help="input directory", required=True, type=str)
parser.add_argument(output_dir_arg, help="output directory", required=True, type=str)
parser.add_argument(files_selector_arg, help="files selector pattern", required=False, type=str)

def _list_files(input_dir : str, selector : str) -> typing.List[str] :
    raw_result = os.listdir(input_dir)
    if selector != "":
        result = list(filter(lambda x: fnmatch.fnmatch(x, selector) if os.path.isfile(os.path.join(input_dir, x)) else True,
                        raw_result))
        return result
    return raw_result

"""
    This function processes files to replace \r\n to \n
    --input and --output are required parameters, --selector is optional, if selector is missing then all files
    must be selected, examples:
    1. select all .sh files from . and save back to . -> lf_fixer.py --input=. --output=. --selector=*.sh
    2. select all files from . and save to the ./out -> lf_fixer.py --input=. --output=./out 
"""
def main() :
    print("Starting to fix line endings CRLF -> LF")
    args = parser.parse_args()
    output_dir = args.output
    input_dir = args.input
    file_selector = args.selector

    files_and_dirs = _list_files(input_dir, file_selector)
    print(files_and_dirs)
    pass

if __name__ == "__main__":
    main()