#include "shell.hh"

using namespace Fract::Shell;
using namespace Fract::Utilities;
using namespace Fract::Objects;

void
shell::printc(std::string msg, std::string color) {
  std::cout << color;
  std::cout << msg;
  std::cout << color::reset;
}

void
shell::print_error(std::string msg) {
  shell::printc(msg + "\n", color::to_ansi(230, 41, 79));
}

std::string
shell::get_input(std::string msg, std::string color) {
  shell::printc(msg, color);
  std::string input;
  std::getline(std::cin, input);
  return string::trim(input);
}

std::string
shell::get_input(std::string msg) {
  std::string input;
  std::cout << msg;
  std::getline(std::cin, input);
  return string::trim(input);
}

std::string
shell::get_input() {
  return shell::get_input(pwd_mark);
}
