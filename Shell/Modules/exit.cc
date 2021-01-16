#include "exit.hh"

using namespace Fract::Shell::Modules;

void exit::process(std::string cmd) {
  if (cmd != "") {
    std::cout << "This module can only be used!" << std::endl;
    return;
  }
  std::exit(0);
}
