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
  for(std::vector<token>::iterator it = tokens->begin(); it < tokens->end(); it++)
  { _result.index++;
    if(it->value == token_lparenthes)
    { first = ++it;
      _result.found = true;
      break;
    }
  }

  /* Skip find close parentheses and result ready steps */
  if(!_result.found) {
    goto result;
  }

  /* Find close parentheses */
  for(std::vector<token>::iterator it = first; it < tokens->end(); it++)
  { if(it->value == token_rparenthes)
    { last = ++it;
      break;
    }
  }

  /* Set range of result */
  _result.range = std::vector<token>(first, last - 1);

  for(int counter = 0; counter <= _result.range.size(); counter++) {
    tokens->erase(last);
  }

  result:
  return _result;
}
