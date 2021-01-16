// Copyright (c) 2021 Fract
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

#include <iostream>

#include "../Include/operating_system.hh"
#include "../Shell/command_processor.hh"
#include "../Shell/shell.hh"
#include "../Shell/Modules/exit.hh"
#include "../Shell/Modules/help.hh"
#include "../Shell/Modules/version.hh"
#include "../Objects/color.hh"

using namespace Fract::Shell;
using namespace Fract::Utilities;

/**
 * @brief Process command.
 * @param ns Namespace of command.
 * @param cmd Command without namespace.
 */
void processCommand(std::string ns, std::string cmd) {
  if (ns == "help") Modules::help::process(cmd);
  else if (ns == "exit") Modules::exit::process(cmd);
  else if (ns == "version") Modules::version::process(cmd);
  else std::cout << "There is no such command!" << std::endl;
}

/**
 * @fn main
 * @brief Entry point
 * @param argc Count of arguments
 * @param argv Arguments
 * @return Exit code
 */
int main(int argc, char const* argv[]) {
  if (argc > 1) { // Started with arguments.
    std::string command = argv[1];
    for(int index = 2; index < argc;)
      command = command + " " + argv[index++];
    processCommand(command_processor::getNamespace(command),
                   command_processor::removeNamespace(command));
    return EXIT_SUCCESS;
  }

  while(true) {
    std::string input = shell::getInput();
    if (input == "")
      continue;
    processCommand(command_processor::getNamespace(input),
                   command_processor::removeNamespace(input));
  }
  return EXIT_SUCCESS;
}
