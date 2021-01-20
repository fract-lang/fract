#include "color.hh"

using namespace Fract::Objects;

std::string color::white = "\033[1;37m";
std::string color::red = "\033[0;31m";
std::string color::green = "\033[0;32m";
std::string color::yellow = "\033[0;33m";
std::string color::blue = "\033[0;34m";
std::string color::magenta = "\033[0;35m";
std::string color::cyan = "\033[0;36m";
std::string color::bold_red = "\033[1;31m";
std::string color::bold_green = "\033[1;32m";
std::string color::bold_yellow = "\033[01;33m";
std::string color::bold_magenta = "\033[1;35m";
std::string color::bold_cyan = "\033[1;36m";
std::string color::reset = "\033[0m";

void
color::enable_vtp() {
#if _WIN32
  HANDLE hOut = GetStdHandle(STD_OUTPUT_HANDLE);
  DWORD dwMode = 0;
  GetConsoleMode(hOut, &dwMode);
  dwMode |= 0x0004;  // Add Virtual Terminal Processing Enable Code.
  SetConsoleMode(hOut, dwMode);
#endif  // _WIN32
}

std::string
color::to_ansi(color color) {
  return color::to_ansi(color.r, color.g, color.b);
}

std::string
color::to_ansi(unsigned short r, unsigned short g, unsigned short b) {
  return "\033[38;2;" + std::to_string(r) + ";" + std::to_string(g) + ";" +
         std::to_string(b) + "m";
}

color::color(unsigned short r, unsigned short g, unsigned short b) {
  this->r = r;
  this->g = g;
  this->b = b;
}
