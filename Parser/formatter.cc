#include "formatter.hh"

using namespace Fract::Parser;

rmlx_range_result
formatter::rmlx_range(std::vector<token> *tokens) {
  rmlx_range_result _result;
  _result.index = -1;
  _result.found = false;

  std::vector<token>::iterator first; /* After iterator of open parentheses */
  std::vector<token>::iterator last; /* After iterator of close parentheses */

  /* Find open parentheses */
  for(int index = 0; index < tokens->size(); index++)
  { std::vector<token>::iterator it = tokens->begin() + index;
    if(it->type == type_open_parenthes) {
      first = ++it;
      _result.index = index;
      _result.found = true;
      break;
    }
  }

  /* Skip find close parentheses and result ready steps */
  if(!_result.found) {
    return _result;
  }

  /* Find close parentheses */
  int count = 1;
  for(int index = _result.index + 1; index < tokens->size(); index++)
  { std::vector<token>::iterator it = tokens->begin() + index;
    if(it->type == type_close_parenthes)
    { count--;
      if(count == 0)
      { last = it;
        break;
      }
    }
    else if(it->type == type_open_parenthes) {
      count++;
    }
    _result.range.push_back(*it);
  }

  /* Remove range from original tokens */
  for(int counter = 0; counter <= _result.range.size(); counter++) {
    tokens->erase(tokens->begin() + _result.index);
  }
  tokens->erase(tokens->begin() + _result.index);

  return _result;
}
