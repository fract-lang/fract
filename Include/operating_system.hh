#ifndef __OPERATING_SYSTEM_HH
#define __OPERATING_SYSTEM_HH

#include <iostream>

#ifdef _WIN32
#include <direct.h>
#else
#include <unistd.h>
#endif  // _WIN32

/// @brief Working directory.
#define pwd getcwd(NULL, 0)

#ifdef _WIN32
/// @brief Path seperator of operating system.
#define path_seperator = '\\'
#else
/// @brief Path seperator of operating system.
#define path_seperator = '/'
#endif

#endif // __OPERATING_SYSTEM_HH
