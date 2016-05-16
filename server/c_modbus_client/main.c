#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <modbus/modbus.h>
#include <time.h>
#include "configuration_parser.h"

void print_config(Configuration *config){
  printf("config->slaves_count: %d\n",config->slaves_count);
  int i;
  int j;
  for(i=0;i<config->slaves_count;i++){
    printf("******slave %d********\n",i+1);
    printf("address: %d\n",config->slaves[i].address);
    printf("description: %s\n",config->slaves[i].description);
    printf("register_count: %d\n",config->slaves[i].registries_count);
    for(j=0;j<config->slaves[i].registries_count;j++){
      printf("***register %d***\n",j+1);
      printf("reg location: %d\n",config->slaves[i].registries[j].location);
      printf("reg length: %d\n",config->slaves[i].registries[j].length);
      printf("reg min: %d\n",config->slaves[i].registries[j].min);
      printf("reg max: %d\n",config->slaves[i].registries[j].max);
      printf("reg type: %d\n",config->slaves[i].registries[j].type);
      printf("reg value_type: %d\n",config->slaves[i].registries[j].value_type);
    }
  }
}

int generate_value(int min, int max){
  if(max<min){
    int temp = min;
    min = max;
    max = temp;
  }
  srand(time(NULL));
  int r = rand();
  r = r % (max-min);

  return r;
}

//TODO implement other types besides holding registers
int start_modbus_client(Configuration *config, char parity, int data_bit, int stop_bit, char *error){
  modbus_t *modbus;

  if(config->slaves_count==0 || config->slaves[0].registries_count==0){
    strcpy(error,"Invalid configuration");
    return 1;
  }

  //we will only use the first slave for now
  int i;
  int start_output_coils_address=0;
  int no_of_output_coils = 0;
  int start_input_coils_address=0;
  int no_of_input_coils = 0;
  int start_holding_registers_address=0;
  int no_of_holding_registers = 0;
  int start_input_registers_address=0;
  int no_of_input_registers = 0;


  for(i=0;i<config->slaves[0].registries_count;i++){
    if(config->slaves[0].registries[i].type == coil){

    }else if(config->slaves[0].registries[i].type == input_discrete){

    }else if(config->slaves[0].registries[i].type == holding){
      no_of_holding_registers++;
      if(no_of_holding_registers==1){
        start_holding_registers_address = config->slaves[0].registries[i].location;
      }
    }else{

    }
  }
  //(output coil, input coil, holding registers,input registers)
  modbus_mapping_t *mapping = modbus_mapping_new_start_address(
    start_output_coils_address,no_of_output_coils,
    start_input_coils_address,no_of_input_coils,
    start_holding_registers_address,no_of_holding_registers,
    start_input_registers_address,no_of_input_registers);
  if(!mapping){
    snprintf(error,ERROR_SIZE,"Failed to allocate the mapping: %s", modbus_strerror(errno));
  }

  modbus = modbus_new_rtu(config->conn.port, config->conn.speed, parity, data_bit, stop_bit);
  if (modbus == NULL) {
    snprintf(error,ERROR_SIZE,"Unable to create the libmodbus context: %s", modbus_strerror(errno));
    return 1;
  }

  if(modbus_set_slave(modbus,3)==-1){
      strcpy(error,"Invalid slave id");
      modbus_free(modbus);
      return 1;
  }

  if(modbus_connect(modbus)==-1){
    snprintf(error,ERROR_SIZE,"Connection failed: %s",modbus_strerror(errno));
    modbus_free(modbus);
    return 1;
  }



  uint8_t req[MODBUS_RTU_MAX_ADU_LENGTH];// request buffer
  int len;// length of the request/response

  while(1) {
    int pos = start_output_coils_address+no_of_output_coils+
      start_input_coils_address+no_of_input_coils;
    mapping->tab_registers[pos] = generate_value(config->slaves[0].registries[0].min,config->slaves[0].registries[0].max);
    len = modbus_receive(modbus, req);
    if (len == -1) break;

    len = modbus_reply(modbus, req, len, mapping);
    if (len == -1) break;
  }

  printf("Exit the loop: %s\n", modbus_strerror(errno));

  modbus_mapping_free(mapping);

  modbus_close(modbus);
  modbus_free(modbus);

  return 0;
}

int main(int argc, char **args){
  char *error;
  FILE *fh = fopen("config.yaml","r");

	Configuration config;
	config.init = init_configuration;
	config.init(&config);

	error = malloc(ERROR_SIZE*sizeof(char));

	if(fh == NULL){
		fprintf(stderr, "file %s could not be opened","config.yaml");
		free(error);
		return 1;
	}
	int rez = parse_yaml_config(fh,&config,error);
	if(rez!=0){
		fprintf(stderr, "%s\n", error);
		free(error);
		return 1;
	}

  rez = start_modbus_client(&config,'N',8,1,error);

  if(rez!=0){
    fprintf(stderr, "%s\n", error);
    return 1;
  }

  //print_config(&config);

  return 0;
}
