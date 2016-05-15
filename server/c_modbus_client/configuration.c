#include "configuration.h"
#include <stdlib.h>

void init_configuration(struct configuration *config){
  config->conn.port[0]='\0';
  config->conn.speed = 0;
  config->slaves_count = 0;
  config->slaves = NULL;
}
