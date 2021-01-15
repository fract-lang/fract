# Copyright (c) 2021 Fract
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# VARIABLES
# GNU GCC C++ compiler.
GCC = g++
# GNU GCC C++ compiler with -c argument for headers.
GCCH = $(GCCH) -c
# GNU GCC output parameter.
OUT = -o

# Name of output.
NAME = Fract

# The "Fract" directory.
DIR_FRACT = $(NAME)
# The "Include" directory.
DIR_INCLUDE = Include
# The "Utilities" directory.
DIR_UTILITIES = Utilities

# Source tree of "Utilities"
define TREE_UTILITIES
$(DIR_UTILITIES)/file_system.o \
$(DIR_UTILITIES)/string.o
endef

# WORKFLOW
# All workflows of this makefile.
all: headers compile
# Headers works.
headers: $(TREE_UTILITIES)

# SUB FLOWS
# All works.
# INCLUDE_UTILITIES
string.o: $(DIR_UTILITIES)/string.cc
	$(GCCH) $< $(OUT) $@

file_system.o: $(DIR_UTILITIES)/file_system.cc
	$(GCCH) $< $(OUT) $@

# Compile the Fract interpreter.
compile: $(DIR_FRACT)/main.cc
	$(GCC) $< $(TREE_UTILITIES) $(OUT) $(NAME)