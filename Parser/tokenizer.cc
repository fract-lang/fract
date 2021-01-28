#include "tokenizer.hh"

using namespace Fract::Parser;
using namespace Fract::Utilities;

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
  std::cout << "TOKENIZER ERROR" << std::endl;
  std::cout << "MESSAGE: " << message << std::endl;
  std::cout << "LINE: " << line_iterator->line << std::endl;
  std::cout << "COLUMN: " << column << std::endl;
  exit(EXIT_FAILURE);
}

/* Last putted token */
token last_token;

token
tokenizer::next_token() {
  token _token;
  _token.line = line_iterator->line;
  _token.column = column;

  /* Return empty token is statement is finished. */
  if(column >= line_iterator->text.length()) {
    return _token;
  }

  std::string statement = line_iterator->text.substr(column - 1);

  /* Ignore whitespaces and tabs. */
  for(int index = 0; index < statement.length(); index++)
  { char ch = statement[index];
    if(
      ch == ' ' ||
      ch == '\t'
      ) {
      column++;
    }
    else {
      statement = statement.substr(index);
      break;
    }
  }

  /* Return empty token is statement is empty. */
  if(statement == "") {
    return _token;
  }

  /* Arithmetic value check */
  if(statement[0] == token_minus[0] || arithmetic::is_numeric(statement[0])) {
    std::string value = statement[0] == token_minus[0] ? token_minus : "";
    if(
      value == "" ||
      (value != "" && (
        last_token.type == type_operator          ||
        last_token.type == type_open_parenthes    ||
        last_token.type == type_close_parenthes
      ))
    )
    { int index = value.length();
      for(; index < statement.size(); index++)
      { char ch = statement[index];
        if(!arithmetic::is_numeric(ch) && ch != token_dot[0]) {
          break;
        }
        value += ch;
      }
      statement = value;
    }
  }

  /* Check anothers. */
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
  }
  else if(string::starts_with(statement, token_plus))
  { _token.type = type_operator;
    _token.value = token_plus;
  }
  else if(string::starts_with(statement, token_minus))
  { _token.type = type_operator;
    _token.value = token_minus;
  }
  else if(string::starts_with(statement, token_star))
  { _token.type = type_operator;
    _token.value = token_star;
  }
  else if(string::starts_with(statement, token_slash))
  { _token.type = type_operator;
    _token.value = token_slash;
  }
  else if(string::starts_with(statement, token_lparenthes))
  { _token.type = type_open_parenthes;
    _token.value = token_lparenthes;
  }
  else if(string::starts_with(statement, token_rparenthes))
  { _token.type = type_close_parenthes;
    _token.value = token_rparenthes;
  }
  else {
    exit_tokenizer_error("What the?: '" + statement + "'");
  }

  column += _token.value.length();
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

  /* Reset to defaults */
  column = 1;
  last_token.type = type_none;
  last_token.value = "";

  token _token;
  while((_token = next_token()).value != "") {
    tokens.push_back(_token);
    last_token = _token;
  }

  if(line_iterator == file->lines.end()) {
    finish = true;
  }

  line_iterator++;
  return tokens;
}
