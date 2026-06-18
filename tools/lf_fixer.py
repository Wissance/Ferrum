import os
import csv
import argparse
import datetime
import typing
import fnmatch

input_dir_arg = "--i"
output_dir_arg = "--o"
files_selector_arg = "--sel"
recursive_arg = "--r"

parser=argparse.ArgumentParser()
parser.add_argument(input_dir_arg, help="input directory", required=True, type=str)
parser.add_argument(output_dir_arg, help="output directory", required=True, type=str)
parser.add_argument(files_selector_arg, help="files selector pattern", required=False, type=str)
parser.add_argument(recursive_arg, help="recursive arg selection pattern", required=False, action='store_true')

class ProcessOptions:
    def __init__(self, input_dir: str, output_dir: str, selector: str, recursive: bool):
        self._input_dir = input_dir
        self._output_dir = output_dir
        self._selector = selector
        self._recursive = recursive

    @property
    def input_dir(self):
        return self._input_dir

    @property
    def output_dir(self):
        return self._output_dir

    @property
    def selector(self):
        return self._selector
    
    @property
    def recursive(self):
        return self._recursive

    _input_dir = ""
    _output_dir = ""
    _selector = ""
    _recursive = False

class FileProcessor:
    def __init__(self, options: ProcessOptions):
        self._options = options

    def list_files(self) -> typing.List[str] :
        raw_result = os.listdir(self._options.input_dir)
        if self._options.selector != "":
            result = list(filter(lambda x: fnmatch.fnmatch(x, self._options.selector) 
                                           if os.path.isfile(os.path.join(self._options.input_dir, x)) 
                                           else self._options.recursive,
                                 raw_result))
            return result
        return raw_result
    
    def process_selected_files(self, files: typing.List[str]):
        for item in files:
            if os.path.isfile(os.path.join(self._options.input_dir, item)):
                pass
            else :
                pass
        pass
    
    _options = None
    
"""
    This function processes files to replace \r\n to \n
    --input and --output are required parameters, --selector is optional, if selector is missing then all files
    must be selected, examples:
    1. select all .sh files from . and save back to . -> lf_fixer.py --input=. --output=. --selector=*.sh
    2. select all files from . and save to the ./out -> lf_fixer.py --input=. --output=./out 
"""
def main() :
    print("###### Starting to fix line endings CRLF -> LF for running scripts in linux ######")
    args = parser.parse_args()
    output_dir = args.o
    input_dir = args.i
    file_selector = args.sel
    recursive_processing = args.r

    options = ProcessOptions(input_dir, output_dir, file_selector, recursive_processing)
    processor = FileProcessor(options)

    files_and_dirs = processor.list_files()
    print(files_and_dirs)
    for item in files_and_dirs:
        if os.path.isfile(os.path.join(input_dir, item)):
            pass
        else :
            pass
    pass
    print("###### Line endings CRLF -> LF for running scripts in linux fixing finished ######")

if __name__ == "__main__":
    main()