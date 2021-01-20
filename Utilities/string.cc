#include "string.hh"

using namespace Fract::Utilities;

std::string
string::trim_start(std::string value) {
  for (int index = 0; index < value.length(); ++index) {
    if (value[index] != ' ') {
      return value.substr(index);
    }
  }
  return value;
}

std::string
string::trim_end(std::string value) {
  for (int index = value.length() - 1; index >= 0; --index) {
    if (value[index] != ' ') {
      return value.substr(0, index + 1);
    }
  }
  return value;
}

std::string
string::trim(std::string value) {
  return trim_end(trim_start(value));
}

std::vector<std::string>
string::split(std::string value,
                                       char seperator) {
  std::vector<std::string> lst;
  int last = 0, index = 0;
  while ((index = value.find(seperator, last)) != std::string::npos) {
    lst.push_back(value.substr(last, index - last));
    last = index + 1;
  }
  if (last != value.length()) {
    lst.push_back(value.substr(last));
  }
  return lst;
}

std::string
string::to_lower(std::string value) {
  for(int index = 0; index < value.length(); index++) {
    value[index] = std::tolower(value[index]);
  }
  return value;
}

std::string
string::to_upper(std::string value) {
  for(int index = 0; index < value.length(); index++) {
    value[index] = std::toupper(value[index]);
  }
  return value;
}

bool
string::starts_with(std::string value, std::string start) {
  if(value.length() < start.length()) {
    return false;
  }
  return value.substr(0, start.length()) == start;
}

bool
string::ends_with(std::string value, std::string end) {
  if(value.length() < end.length()) {
    return false;
  }
  return value.substr(value.length() - end.length()) == end;
}

bool
string::contains(std::string value, std::string check) {
  return value.find(check) != std::string::npos;
}
