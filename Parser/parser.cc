#include "parser.hh"

using namespace Fract::Parser;
using namespace Fract::Utilities;

code_file
parser::ready_file(std::string path) {
  code_file file;
  file.lines = parser::ready_lines(file_system::get_lines_of_file(path));
  file.path = path;
  file.stream = std::ifstream(path);
  return file;
}

std::vector<code_line>
parser::ready_lines(std::vector<std::string> lines) {
  std::vector<code_line> ready_lines;
  for(int index = 0; index < lines.size(); index++) {
    ready_lines.push_back(code_line{index + 1, lines[index]});
  }
  return ready_lines;
}

parser::parser(std::string path, int type) {
  file = parser::ready_file(path);
  _tokenizer = tokenizer(&file);
  this->type = type;
}

void
parser::parse() {
  while(!_tokenizer.finish) {
    std::cout << "-----------------------" << std::endl;
    for(auto _token : _tokenizer.tokenize_next()) {
      std::cout << _token.type << " - "  << _token.value << std::endl;
    }
  }
}
