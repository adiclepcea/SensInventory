#include "configuration.h"
#include <yaml.h>
#include <stdio.h>

//parse_yaml_config - will return 0 for no error and will fill the configuration
//with the parsed values
//will return non zero on failure and will fill the error string with the details
int parse_yaml_config(FILE *f,Configuration *config, char *error);
int parse_by_event(yaml_parser_t parser,Configuration *config,char *error);
