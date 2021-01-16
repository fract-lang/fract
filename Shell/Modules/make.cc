#include "make.hh"

using namespace Fract::Shell::Modules;
using namespace Fract::Utilities;

void make::process(std::string cmd) {
  if(cmd == "")
  { std::cout << "This module cannot only be used!" << std::endl;
    return;
  }
  cmd += !string::endsWith(cmd, fract_extension) ? fract_extension : "";
  if(!file_system::existFile(cmd))
  { std::cout << "The Fract file is not exists: " << cmd << std::endl;
    return;
  }
}

bool make::check(std::string value) {
  if(string::endsWith(value, fract_extension)) return true;
  else value += fract_extension;
  if(file_system::existFile(value)) return true;
  return false;
}
