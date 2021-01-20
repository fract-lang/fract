#include "make.hh"

using namespace Fract::Shell::Modules;
using namespace Fract::Utilities;
using namespace Fract::Parser;

void
make::process(std::string cmd) {
  if(cmd == "")
  { std::cout << "This module cannot only be used!" << std::endl;
    return;
  }
  cmd += !string::ends_with(cmd, fract_extension) ? fract_extension : "";
  if(!file_system::exist_file(cmd))
  { std::cout << "The Fract file is not exists: " << cmd << std::endl;
    return;
  }
  parser entry(cmd, type_entry_file);
}

bool
make::check(std::string value) {
  if(string::ends_with(value, fract_extension)){
    return true;
  }
  else {
    value += fract_extension;
  }
  
  if(file_system::exist_file(value)) {
    return true;
  }
  return false;
}
