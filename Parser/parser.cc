#include "parser.hh"

using namespace Fract::Parser;
using namespace Fract::Utilities;

code_file
parser::readyFile(std::string path) {
  code_file file;
  file.lines = parser::readyLines(file_system::getLinesOfFile(path));
  file.path = path;
  file.stream = std::ifstream(path);
  return file;
}

std::vector<code_line>
parser::readyLines(std::vector<std::string> lines) {
  std::vector<code_line> readyLines;
  for(int index = 0; index < lines.size(); index++) {
    readyLines.push_back(code_line{index + 1, lines[index]});
  }
  return readyLines;
}

parser::parser(std::string path, int type) {
  file = parser::readyFile(path);
  this->type = type;
}
