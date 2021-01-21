#include "tokenizer.hh"

using namespace Fract::Parser;

tokenizer::tokenizer() { /**/ }

tokenizer::tokenizer(code_file *file) {
  finish = false;
  this->file = file;
  column = 1;
  line_iterator = file->lines.begin();
}

void
tokenizer::exit_tokenizer_error(std::string message) {
  std::cout << std::endl;
  std::cout << "ERROR" << std::endl;
  std::cout << "MESSAGE: " << message << std::endl;
  std::cout << "LINE: " << line_iterator->line << std::endl;
  std::cout << "COLUMN: " << column << std::endl;
  exit(EXIT_FAILURE);
}

token
tokenizer::next_token() {
  token _token;
  _token.line = line_iterator->line;
  _token.column = column;

  for(int index = 0; index < line_iterator->text.length(); index++)
  { if(line_iterator->text[index] == ' ') {
      column++;
    }
  }

  if(column == line_iterator->text.length()) {
    return _token;
  }

  std::string statement = line_iterator->text.substr(column - 1);

  if(arithmetic::is_integer_number(statement))
  { _token.type = type_value;
    _token.value = statement;
  }
  else if(arithmetic::is_floating_number(statement))
  { _token.type = type_value;
    _token.value = statement;
  }
  else if(string::starts_with(statement, kw_variable))
  { _token.type = type_let;
    _token.value = kw_variable;
  } else {
    exit_tokenizer_error("What the?: '" + line_iterator->text + "'");
  }

  column += _token.value.length() - 1;

  return _token;
}

std::vector<token>
tokenizer::tokenize_next() {
  std::vector<token> tokens;

  if(finish) {
    return tokens;
  }

  if(line_iterator->text == "") {
    return tokens;
  }

  column = 1;
  token _token;
  while((_token = next_token()).value != "") {
    tokens.push_back(_token);
  }

  if(line_iterator == file->lines.end()) {
    finish = true;
  }

  line_iterator++;
  return tokens;
}
