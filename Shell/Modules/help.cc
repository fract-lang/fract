#include "help.hh"

using namespace Fract::Shell::Modules;
using namespace Fract::Utilities;

void help::process(std::string cmd) {
  if(cmd != "") {
    std::cout << "This module can only be used!" << std::endl;
    return;
  }
  cli::printMapAsTable({
    {"make", "Interprete Fract code."},
    {"version", "Show version."},
    {"help", "Show help."},
    {"exit", "Exit."}
  });
}
