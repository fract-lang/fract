#include "arithmetic.hh"

using namespace Fract;
using namespace Fract::Grammar;

bool
arithmetic::is_integer_number(std::string value) {
  if(value == "") return false;
  else if (value == ".") return false;
	
  for(int index = 0; index < value.length(); index++)
  { if(!arithmetic::is_number(value[index])) {
      return false;
    }
  }
  return true;
}

bool
arithmetic::is_floating_number(std::string value) {
  if(value == "") return false;
  else if (value == ".") return false;

  bool dotted = false;
  for(int index = 0; index < value.length(); index++)
  { char ch = value[index];
    if(ch == token_dot[0] && !dotted)
    { dotted = true;
      continue;
    } else if (ch == token_dot[0] && dotted) {
      return false;
    }
    if(!arithmetic::is_number(value[index])) {
      return false;
    }
  }
  return true;
}

bool
arithmetic::is_number(char ch) {
  return ch == '0' ||
         ch == '1' ||
         ch == '2' ||
         ch == '3' ||
         ch == '4' ||
         ch == '5' ||
         ch == '6' ||
         ch == '7' ||
         ch == '8' ||
         ch == '9';
}
