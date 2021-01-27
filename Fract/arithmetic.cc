#include "arithmetic.hh"

using namespace Fract;

bool
arithmetic::is_types_compatible(int type0, int type1) {
  return is_integer_type(type0) ?
           is_integer_type(type1) :
           is_float_type(type0) ?
             is_float_type(type1) : false;
}

bool
arithmetic::is_integer_type(int type) {
  return type == type_unsigned_int16 ||
         type == type_unsigned_int32 ||
         type == type_unsigned_int64 ||
         type == type_int16          ||
         type == type_int32          ||
         type == type_int64;
}

bool
arithmetic::is_float_type(int type) {
  return type == type_float ||
         type == type_double;
}

unsigned short
arithmetic::to_uint16(std::string value) {
  unsigned short x = integer_default;
  sscanf(value.c_str(), "%hu", &x);
  return x;
}

unsigned int
arithmetic::to_uint32(std::string value) {
  unsigned int x = integer_default;
  sscanf(value.c_str(), "%u", &x);
  return x;
}

unsigned long
arithmetic::to_uint64(std::string value) {
  unsigned long x = integer_default;
  sscanf(value.c_str(), "%lu", &x);
  return x;
}

short
arithmetic::to_int16(std::string value) {
  short x = integer_default;
  sscanf(value.c_str(), "%hd", &x);
  return x;
}

int
arithmetic::to_int32(std::string value) {
  int x = integer_default;
  sscanf(value.c_str(), "%d", &x);
  return x;
}

long
arithmetic::to_int64(std::string value) {
  long x = integer_default;
  sscanf(value.c_str(), "%ld", &x);
  return x;
}

float
arithmetic::to_float(std::string value) {
  float x = float_default;
  sscanf(value.c_str(), "%f", &x);
  return x;
}

double
arithmetic::to_double(std::string value) {
  double x = float_default;
  sscanf(value.c_str(), "%lf", &x);
  return x;
}

bool
arithmetic::is_integer_number(std::string value) {
  if(value == "") {
    return false;
  }
  else if (value == ".") {
    return false;
  }

  for(int index = value[0] == token_minus[0] ? 1 : 0; index < value.length(); index++)
  { if(!arithmetic::is_numeric(value[index])) {
      return false;
    }
  }
  return true;
}

bool
arithmetic::is_floating_number(std::string value) {
  if(value == "") {
    return false;
  }
  else if (value == ".") {
    return false;
  }

  bool dotted = false;
  for(int index = value[0] == token_minus[0] ? 1 : 0; index < value.length(); index++)
  { char ch = value[index];
    if(ch == token_dot[0] && !dotted)
    { dotted = true;
      continue;
    } else if (ch == token_dot[0] && dotted) {
      return false;
    }
    if(!arithmetic::is_numeric(value[index])) {
      return false;
    }
  }
  return true;
}

bool
arithmetic::is_numeric(char ch) {
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

unsigned short
arithmetic::to_numeric(char ch) {
  return ch == '0' ? 0 :
         ch == '1' ? 1 :
         ch == '2' ? 2 :
         ch == '3' ? 3 :
         ch == '4' ? 4 :
         ch == '5' ? 5 :
         ch == '6' ? 6 :
         ch == '7' ? 7 :
         ch == '8' ? 8 : 9;
}

bool
arithmetic::bigger(char one, char two) {
  return to_numeric(one) > to_numeric(two);
}

bool
arithmetic::lower(char one, char two) {
  return to_numeric(one) < to_numeric(two);
}

bool
arithmetic::bigger_str(std::string one, std::string two) {
  bool first_is_floating =
    one.find(token_dot) != std::string::npos;
  bool second_is_floating =
    two.find(token_dot) != std::string::npos;

  std::string first = one;
  std::string second = two;

  if(first_is_floating) {
    first = one.substr(0, one.find(token_dot));
  }
  if(second_is_floating) {
    second = two.substr(0, two.find(token_dot));
  }

  if(first.length() > second.length()) {
    return true;
  }

  if(first == second)
  { if(first_is_floating || second_is_floating) {
      goto float_check;
    }
    return false;
  }

  for(int index = 0; index < second.length(); index++)
  { if(bigger(first[index], second[index])) {
      return true;
    }
    else if(lower(first[index], second[index])) {
      return false;
    }
  }

  if(!first_is_floating && !second_is_floating) {
    return false;
  }

  // *************
  //  FLOAT CHECK
  // *************

  float_check:
  first =
    first_is_floating ? one.substr(one.find(token_dot) + 1) : "";
  second =
    second_is_floating ? two.substr(two.find(token_dot) + 1) : "";

  std::size_t first_len = first.length();
  std::size_t second_len = second.length();

  if(first_is_floating && !second_is_floating)
  { for(int index = 0; index < first_len; index++)
    { if(first[index] != '0') {
        return true;
      }
    }
  }
  else if(!first_is_floating && second_is_floating)
  { for(int index = 0; index < second_len; index++)
    { if(second[index] != '0') {
        return false;
      }
    }
  }

  if(first_len != second_len)
  { if(first_len > second_len)
    { for(int index = second_len; index < first_len; index++)
      { if(first[index] != '0') {
          return true;
        }
      }
    }
    else
    { for(int index = first_len; index < second_len; index++)
      { if(second[index] != '0') {
          return false;
        }
      }
    }
  }

  for(int index = 0; index < first_len; index++)
  { if(bigger(first[index], second[index])) {
      return true;
    }
  }

  return false;
}

bool
arithmetic::lower_str(std::string one, std::string two) {
  bool first_is_floating =
    one.find(token_dot) != std::string::npos;
  bool second_is_floating =
    two.find(token_dot) != std::string::npos;

  std::string first = one;
  std::string second = two;

  if(first_is_floating) {
    first = one.substr(0, one.find(token_dot));
  }
  if(second_is_floating) {
    second = two.substr(0, two.find(token_dot));
  }

  if(first.length() < second.length()) {
    return true;
  }

  if(first == second)
  { if(first_is_floating || second_is_floating) {
      goto float_check;
    }
    return false;
  }

  for(int index = 0; index < second.length(); index++)
  { if(lower(first[index], second[index])) {
      return true;
    }
    else if(bigger(first[index], second[index])) {
      return false;
    }
  }

  if(!first_is_floating && !second_is_floating) {
    return false;
  }

  // *************
  //  FLOAT CHECK
  // *************

  float_check:
  first =
    first_is_floating ? one.substr(one.find(token_dot) + 1) : "";
  second =
    second_is_floating ? two.substr(two.find(token_dot) + 1) : "";

  std::size_t first_len = first.length();
  std::size_t second_len = second.length();

  if(first_is_floating && !second_is_floating)
  { for(int index = 0; index < first_len; index++)
    { if(first[index] != '0') {
        return false;
      }
    }
  }
  else if(!first_is_floating && second_is_floating)
  { for(int index = 0; index < second_len; index++)
    { if(second[index] != '0') {
        return true;
      }
    }
  }

  if(first_len != second_len)
  { if(first_len > second_len)
    { for(int index = second_len; index < first_len; index++)
      { if(first[index] != '0') {
          return false;
        }
      }
    }
    else
    { for(int index = first_len; index < second_len; index++)
      { if(second[index] != '0') {
          return true;
        }
      }
    }
  }

  for(int index = 0; index < first_len; index++)
  { if(lower(first[index], second[index])) {
      return true;
    }
  }

  return false;
}

bool
arithmetic::equals_str(std::string one, std::string two) {
  bool first_is_floating =
    one.find(token_dot) != std::string::npos;
  bool second_is_floating =
    two.find(token_dot) != std::string::npos;

  std::string first = one;
  std::string second = two;

  if(first_is_floating) {
    first = one.substr(0, one.find(token_dot));
  }
  if(second_is_floating) {
    second = two.substr(0, two.find(token_dot));
  }

  if(!first_is_floating && !second_is_floating) {
    return first == second;
  }
  else
  { if(first != second) {
      return false;
    }
  }

  // *************
  //  FLOAT CHECK
  // *************

  first =
    first_is_floating ? one.substr(one.find(token_dot) + 1) : "";
  second =
    second_is_floating ? two.substr(two.find(token_dot) + 1) : "";

  std::size_t first_len = first.length();
  std::size_t second_len = second.length();

  if(first_is_floating && !second_is_floating)
  { for(int index = 0; index < first_len; index++)
    { if(first[index] != '0') {
        return false;
      }
    }
    return true;
  }
  else if(!first_is_floating && second_is_floating)
  { for(int index = 0; index < second_len; index++)
    { if(second[index] != '0') {
        return false;
      }
    }
    return true;
  }

  if(first_len != second_len)
  { if(first_len > second_len)
    { for(int index = second_len; index < first_len; index++)
      { if(first[index] != '0') {
          return false;
        }
      }
    }
    else
    { for(int index = first_len; index < second_len; index++)
      { if(second[index] != '0') {
          return false;
        }
      }
    }
  }

  for(int index = first_len; index < first_len; index++)
  { if(first[index] != second[index]) {
      return false;
    }
  }

  return true;
}
