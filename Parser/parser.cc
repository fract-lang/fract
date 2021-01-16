#include "parser.hh"

using namespace Fract::Parser;

std::vector<code_line> parser::readyLines(std::vector<std::string> lines) {
  std::vector<code_line> readyLines;
  for(int index = 0; index < lines.size(); index++)
    readyLines.push_back(code_line{index + 1, lines[index]});
  return readyLines;
}
