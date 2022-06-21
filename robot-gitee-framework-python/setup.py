# coding: utf-8

from setuptools import setup, find_packages

NAME = "robotgitee"
VERSION = "1.0.0"
# To install the library, run the following
#
# python3 setup.py bdist_wheel
# pip install ./dist/***.whl
#
# prerequisite: setuptools
# http://pypi.python.org/pypi/setuptools

REQUIRES = [
    "gitee>=1.0.0",
]

setup(
    name=NAME,
    version=VERSION,
    description="robot gitee framework",
    author_email="",
    url="",
    keywords=["robot gitee framework"],
    install_requires=REQUIRES,
    packages=find_packages(),
    include_package_data=True,
)
