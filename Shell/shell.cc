#include "shell.hh"

using namespace Fract::Shell;
using namespace Fract::Utilities;
using namespace Fract::Objects;

void shell::printc(std::string msg, std::string color) {
  std::cout << color;
  std::cout << msg;
  std::cout << color::reset;
}

void shell::printError(std::string msg) {
  shell::printc(msg + "\n", color::toANSI(230, 41, 79));
}

std::string shell::getInput(std::string msg, std::string color) {
  shell::printc(msg, color);
  std::string input;
  std::getline(std::cin, input);
  return string::trim(input);
}

std::string shell::getInput(std::string msg) {
  std::string input;
  std::cout << msg;
  std::getline(std::cin, input);
  return string::trim(input);
}

std::string shell::getInput() {
  return shell::getInput(pwd_mark);
}
