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
  std::cout << "TOKENIZER ERROR" << std::endl;
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

  /* If numeric? */
  if(arithmetic::is_numeric(statement[0]))
  { for(int index = 0; index < statement.length(); index++)
    { char ch = statement[index];
      if(ch == token_dot[0]) {
        continue;
      }
      else if(!arithmetic::is_numeric(ch))
      { statement = statement.substr(0, index);
        break;
      }
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
    column++;
  }
  else if(string::starts_with(statement, token_minus))
  { _token.type = type_operator;
    _token.value = token_minus;
    column++;
  }
  else if(string::starts_with(statement, token_star))
  { _token.type = type_operator;
    _token.value = token_star;
    column++;
  }
  else if(string::starts_with(statement, token_slash))
  { _token.type = type_operator;
    _token.value = token_slash;
    column++;
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
