import os
import csv
import argparse
import datetime
import typing

input_dir_arg = "--input"
output_dir_arg = "--output"
files_selector_arg = "--selector"

parser=argparse.ArgumentParser()
parser.add_argument(input_dir_arg, help="input directory", required=True, type=str)
parser.add_argument(output_dir_arg, help="output directory", required=True, type=str)
parser.add_argument(files_selector_arg, help="files selector pattern", required=False, type=str)

"""
    This function processes files to replace \r\n to \n
    --input and --output are required parameters, --selector is optional, if selector is missing then all files
    must be selected, examples:
    1. select all .sh files from . and save back to . -> lf_fixer.py --input=. --output=. --selector=*.sh
    2. select all files from . and save to the ./out -> lf_fixer.py --input=. --output=./out 
"""
def main() :
    pass