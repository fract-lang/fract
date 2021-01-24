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
#include "../Shell/Modules/make.hh"
#include "../Shell/Modules/version.hh"
#include "../Objects/color.hh"
#include "../Utilities/file_system.hh"
#include "arithmetic.hh"

using namespace Fract::Shell;
using namespace Fract::Utilities;

/**
 * @brief Process command.
 * @param ns Namespace of command.
 * @param cmd Command without namespace.
 */
void
process_command(std::string ns, std::string cmd) {
  if (ns == "help") {
    Modules::help::process(cmd);
  }
  else if (ns == "exit") {
    Modules::exit::process(cmd);
  }
  else if (ns == "version") {
    Modules::version::process(cmd);
  }
  else if (ns == "make") {
    Modules::make::process(cmd);
  }
  else if (Modules::make::check(ns)) {
    Modules::make::process(ns + cmd);
  }
  else {
    std::cout << "There is no such command!" << std::endl;
  }
}

/**
 * @fn main
 * @brief Entry point
 * @param argc Count of arguments
 * @param argv Arguments
 * @return Exit code
 */
int
main(int argc, char const* argv[]) {
  /*while(true)
  { std::string x;
    std::string y;
    std::cin >> x;
    std::cin >> y;
    std::cout << Fract::arithmetic::lower_str(x, y) << std::endl;
  }*/
  /*while(true)
  { std::string x;
    std::cin >> x;
    std::cout << Fract::arithmetic::to_double(x) << std::endl;
  }*/

  if (argc <= 1) {// Not started with arguments.
    return EXIT_SUCCESS;
  }

  std::string command = argv[1];

  for(int index = 2; index < argc;) {
    command = command + " " + argv[index++];
  }

  process_command(command_processor::get_namespace(command),
                 command_processor::remove_namespace(command));

  return EXIT_SUCCESS;
}
