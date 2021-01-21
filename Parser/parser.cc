#include "parser.hh"

using namespace Fract::Parser;
using namespace Fract::Utilities;

code_file
parser::ready_file(std::string path) {
  code_file file;
  file.lines = parser::ready_lines(
    file_system::get_lines_of_file(path));
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

void
parser::exit_parser_error(token _token, std::string message) {
  std::cout << std::endl;
  std::cout << "ERROR" << std::endl;
  std::cout << "MESSAGE: " << message << std::endl;
  std::cout << "LINE: " << _token.line << std::endl;
  std::cout << "COLUMN: " << _token.column << std::endl;
  exit(EXIT_FAILURE);
}

parser::parser(std::string path, int type) {
  file = parser::ready_file(path);
  _tokenizer = tokenizer(&file);
  this->type = type;
}

void
parser::parse() {
  while(!_tokenizer.finish) {
    std::vector<token> tokens = _tokenizer.tokenize_next();
    std::vector<token>::iterator first = tokens.begin();
    if(first->type == type_value) {
      print_value(process_value(tokens, &first));
    }
    else {
      exit_parser_error(*first, "What the?:" + first->value);
    }
    /*std::cout << "-----------------------" << std::endl;
    for(auto _token : _tokenizer.tokenize_next()) {
      std::cout << _token.type << " - "  << _token.value << std::endl;
    }*/
  }
}

void
parser::print_value(value _value) {
  if(
    _value.type == type_int16             ||
    _value.type == type_int32             ||
    _value.type == type_int64             ||
    _value.type == type_int64_64          ||
    _value.type == type_byte              ||
    _value.type == type_signed_byte       ||
    _value.type == type_unsigned_int16    ||
    _value.type == type_unsigned_int32    ||
    _value.type == type_unsigned_int32    ||
    _value.type == type_unsigned_int64    ||
    _value.type == type_unsigned_int64_64 ||
    _value.type == type_float             ||
    _value.type == type_double
    ) {
    std::cout << _value.content << std::endl;
  }
  else {
    std::cout << _value.content << std::endl;
  }
}

value
parser::process_value(std::vector<token> tokens,
                      std::vector<token>::iterator *it) {
  value _value;
  if(arithmetic::is_integer_number((*it)->value))
  { _value.content = (*it)->value;
    _value.type = type_int32;
  }
  else if(arithmetic::is_floating_number((*it)->value))
  { _value.content = (*it)->value;
    _value.type = type_double;
  }
  return _value;
}
