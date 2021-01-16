#include "version.hh"

using namespace Fract::Shell::Modules;

void version::process(std::string cmd) {
  if (cmd != "") {
    std::cout << "This module can only be used!" << std::endl;
    return;
  }
  std::cout << "Fract Version [" fract_version "]" << std::endl;
}
