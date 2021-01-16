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
NAME = fract

# The "Fract" directory.
DIR_FRACT = Fract
# The "Include" directory.
DIR_INCLUDE = Include
# The "Objects" directory.
DIR_OBJECTS = Objects
# The "Shell" directory.
DIR_SHELL = Shell
# The "Shell/Modules" directory.
DIR_SHELL_MODULES = $(DIR_SHELL)/Modules
# The "Utilities" directory.
DIR_UTILITIES = Utilities

# Include tree of "Objects"
define TREE_OBJECTS
$(DIR_OBJECTS)/color.cc
endef

# Include tree of "Shell"
define TREE_SHELL
$(DIR_SHELL)/command_processor.cc \
$(DIR_SHELL)/shell.cc
endef

# Include tree of "Shell/Modules"
define TREE_SHELL_MODULES
$(DIR_SHELL_MODULES)/exit.cc \
$(DIR_SHELL_MODULES)/help.cc \
$(DIR_SHELL_MODULES)/make.cc \
$(DIR_SHELL_MODULES)/version.cc
endef

# Include tree of "Utilities"
define TREE_UTILITIES
$(DIR_UTILITIES)/cli.cc \
$(DIR_UTILITIES)/file_system.cc \
$(DIR_UTILITIES)/string.cc
endef

# WORKFLOW
# All workflows of this makefile.
all: compile

# Compile the Fract interpreter.
compile: $(DIR_FRACT)/main.cc
	$(GCC) $< $(TREE_OBJECTS) $(TREE_SHELL) $(TREE_SHELL_MODULES) \
	$(TREE_UTILITIES) $(OUT) $(NAME)