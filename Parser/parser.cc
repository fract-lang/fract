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
    ready_lines.push_back(code_line{index + 1, lines[index] + " "});
  }
  return ready_lines;
}

void
parser::exit_parser_error(token _token, std::string message) {
  std::cout << std::endl;
  std::cout << "PARSER ERROR" << std::endl;
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
  while(!_tokenizer.finish)
  { std::vector<token> tokens = _tokenizer.tokenize_next();
    std::vector<token>::iterator first = tokens.begin();
    check_parentheses(&tokens);
    if(first->type == type_value) {
      print_value(process_value(&tokens, &first));
    }
    else {
      exit_parser_error(*first, "What the?:" + first->value);
    }
  }
}

void
parser::print_value(value _value) {
  if(
    _value.type == type_int16             ||
    _value.type == type_int32             ||
    _value.type == type_int64             ||
    _value.type == type_byte              ||
    _value.type == type_signed_byte       ||
    _value.type == type_unsigned_int16    ||
    _value.type == type_unsigned_int32    ||
    _value.type == type_unsigned_int32    ||
    _value.type == type_unsigned_int64    ||
    _value.type == type_float
    ) {
    std::cout << _value.content << std::endl;
  }
  else {
    std::cout << _value.content << std::endl;
  }
}

value
parser::process_value(std::vector<token> *tokens,
                      std::vector<token>::iterator *it) {
  value _value;
  int type = ptype_none;
  for(; *it < (*tokens).end(); (*it)++)
  { std::string _cache_value = _value.content;
    int _cache_type = _value.type;

    /* Check operators. */
    if((*it)->value == token_plus)
    { type = ptype_addition;
      continue;
    }
    else if((*it)->value == token_minus)
    { type = ptype_subtraction;
      continue;
    }
    else if((*it)->value == token_star)
    { type = ptype_multiplication;
      continue;
    }
    else if((*it)->value == token_slash)
    { type = ptype_division;
      continue;
    }

    if(arithmetic::is_integer_number((*it)->value))
    { _value.content = (*it)->value;
      _value.type = type_int32;
    }
    else if(arithmetic::is_floating_number((*it)->value))
    { _value.content = (*it)->value;
      _value.type = type_float;
    }
    else {
      exit_parser_error(**it, "What the?: " + (*it)->value);
    }

    /* If not exists any operator. */
    if(type == ptype_none) {
      continue;
    }

    /* If data types are not compatible! */
    if(!arithmetic::is_types_compatible(_cache_type, _value.type)) {
      exit_parser_error(**it, "Data types is not compatible!");
    }

    double _arithmetic_value = arithmetic::to_double(_cache_value);
    double _cache_arithmetic_value = arithmetic::to_double(_value.content);

    if(type == ptype_addition) {
      _value.content =
        std::to_string(_arithmetic_value + _cache_arithmetic_value);
    }
    else if(type == ptype_subtraction) {
      _value.content =
        std::to_string(_arithmetic_value - _cache_arithmetic_value);
    }
    else if(type == ptype_division) {
      _value.content =
        std::to_string(_arithmetic_value / _cache_arithmetic_value);
    }
    else if(type == ptype_multiplication) {
      _value.content =
        std::to_string(_arithmetic_value * _cache_arithmetic_value);
    }

    /* Reset type. */
    type = ptype_none;
  }

  /* If exists unprocessed operator? */
  if(type != ptype_none) {
    exit_parser_error(*((*it) - 1), "Unused operator?");
  }

  return _value;
}

void
parser::check_parentheses(std::vector<token> *tokens) {
  int count = 0;
  std::vector<token>::iterator lastOpen;
  for(std::vector<token>::iterator it = tokens->begin(); it < tokens->end(); it++)
  { if((*it).type == type_open_parenthes) {
      lastOpen = it;
      count++;
    }
    else if((*it).type == type_close_parenthes)
    { if(count == 0) {
        exit_parser_error(*it, "The extra parentheses are closed!");
      }
      count--;
    }
  }
  if(count > 0) {
    exit_parser_error(*lastOpen, "The parentheses are opened but not closed!");
  }
}
