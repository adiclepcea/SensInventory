#include "configuration_parser.h"
#include "termios.h"

int parse_yaml_config(FILE *f,Configuration *config, char *error){
  yaml_parser_t parser;
	yaml_token_t token;
	if (!yaml_parser_initialize(&parser)){
    snprintf(error, ERROR_SIZE,"%s","failed to initialize the parser");
		return 1;
	}
  yaml_parser_set_input_file(&parser, f);

  int rez = parse_by_event(parser,config, error);

  yaml_parser_delete(&parser);

  return rez;
}

int get_speed(unsigned long speed){
  int rez = 0;
  switch(speed){
    case 115200:
      rez = B115200;
      break;
    case 57600:
      rez = B57600;
      break;
    case 38400:
      rez = B38400;
      break;
    case 19200:
      rez = B19200;
      break;
    case 9600:
      rez = B9600;
      break;
    case 4800:
      rez = B4800;
      break;
    case 2400:
      rez = B2400;
      break;
    case 1200:
      rez = B1200;
      break;
    default:
      rez = 0;
      break;
  }
  return rez;
}
//
int get_next_scalar_event(yaml_parser_t *parser,yaml_event_t *event, char *error, char *section, char *message){
  if(!yaml_parser_parse(parser, event)){
    snprintf(error,ERROR_SIZE,"Error parsing yaml. %s section",section);
    yaml_event_delete(event);
    return 1;
  }
  if (event->type != YAML_SCALAR_EVENT){
    snprintf(error,ERROR_SIZE,"Unexpected yaml parse. %s.",message);
    yaml_event_delete(event);
    return 1;
  }
  return 0;
}

void add_registry(Configuration *config){
  config->slaves[config->slaves_count-1].registries_count++;
  config->slaves[config->slaves_count-1].registries = realloc(config->slaves[config->slaves_count-1].registries,config->slaves[config->slaves_count-1].registries_count);
}

int get_int_value(yaml_parser_t *parser,yaml_event_t *event,char *error, char *section, char *message, int *value){
  if(get_next_scalar_event(parser,event,error,"Slave","Expected length")){
    return 1;
  }
  int val = atoi(event->data.scalar.value);
  if(!val && strcmp(event->data.scalar.value,"0")!=0){
    snprintf(error,ERROR_SIZE,"%s is not a valid integer!",event->data.scalar.value);
    return 1;
  }
  *value = val;
  return 0;
}

int parse_registries(yaml_parser_t *parser,Configuration *config, char *error){
  yaml_event_t event;
  //we need to get read of the mapping_start_event

  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Registry section");
    return 1;
  }
  yaml_event_delete(&event);

  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Registry section");
    return 1;
  }

  if(event.type != YAML_MAPPING_START_EVENT){
    snprintf(error,ERROR_SIZE,"Invalid event. Expected MAPPING START, got %d",event.type);
    return 1;
  }
  add_registry(config);
  while(1 && event.type !=YAML_SEQUENCE_END_EVENT && event.type != YAML_NO_EVENT){
    yaml_event_delete(&event);
    if(!yaml_parser_parse(parser, &event)){
      strcpy(error,"Error parsing yaml. Registry section");
      return 1;
    }
    if(event.type==YAML_SCALAR_EVENT){
      if(strcmp(event.data.scalar.value,"length")==0){
        if(get_int_value(parser,&event,error,"registries","length",&config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].length)!=0){
          return 1;
        }
      }else if(strcmp(event.data.scalar.value,"location")==0){
        if(get_int_value(parser,&event,error,"registries","location",&config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].location)!=0){
          return 1;
        }
      }else if(strcmp(event.data.scalar.value,"max")==0){
        if(get_int_value(parser,&event,error,"registries","max",&config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].max)!=0){
          return 1;
        }
      }else if(strcmp(event.data.scalar.value,"min")==0){
        if(get_int_value(parser,&event,error,"registries","min",&config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].min)!=0){
          return 1;
        }
      }else if(strcmp(event.data.scalar.value,"type")==0){
        if(get_next_scalar_event(parser,&event,error,"registries","Expected type")){
          return 1;
        }
        char cIn[21];
        strncpy(cIn,event.data.scalar.value,20);
        if(strcmp(cIn,"holding")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].type=holding;
        }else if(strcmp(cIn,"input")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].type=input;
        }else if(strcmp(cIn,"coil")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].type=coil;
        }else if(strcmp(cIn,"input_discrete")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].type=input_discrete;
        }
      }else if(strcmp(event.data.scalar.value,"value")==0){
        if(get_next_scalar_event(parser,&event,error,"registries","Expected value")){
          return 1;
        }
        char cIn[21];
        strncpy(cIn,event.data.scalar.value,20);
        if(strcmp(cIn,"random_generated")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].value_type=random_generated_value;
        }else if(strcmp(cIn,"fixed")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].value_type=fixed_value;
        }else if(strcmp(cIn,"read")==0){
          config->slaves[config->slaves_count-1].registries[config->slaves[config->slaves_count-1].registries_count-1].value_type=read_value;
        }
      }
    }
    if(event.type==YAML_MAPPING_END_EVENT){ //finished with one register. see if another one folows
      if(!yaml_parser_parse(parser, &event)){
        strcpy(error,"Error parsing yaml. Registry section");
        return 1;
      }
      if(event.type==YAML_MAPPING_START_EVENT){
        add_registry(config);
      }
    }
  }

  return 0;

}

int parse_slave(yaml_parser_t *parser, Configuration *config, char *error){
  yaml_event_t event;
  //we need to get read of the mapping_start_event
  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Slave section");
    return 1;
  }

  //we need to add to the slaves collection one
  config->slaves = realloc(config->slaves,(config->slaves_count+1)*sizeof(Slave));
  if(!config->slaves){
    strcpy(error, "failed to allocate memory for a nea slave");
    return 1;
  }
  config->slaves_count++;
  config->slaves[config->slaves_count-1].registries = NULL;
  config->slaves[config->slaves_count-1].registries_count = 0;

  while(1 && event.type!=YAML_MAPPING_END_EVENT && event.type!=YAML_NO_EVENT){
    //we should have either the address, the description or the receptions scalars
    if(event.type!=YAML_SCALAR_EVENT){
      snprintf(error, ERROR_SIZE,"Error parsing yaml. Slave section. Expecting scalar. Got %d",event.type);
      return 1;
    }
    if(strcmp("address",event.data.scalar.value)==0){
      if(get_next_scalar_event(parser,&event,error,"Slave","Expected address")){
        return 1;
      }
      int address = atoi(event.data.scalar.value);
      if(!address){
        snprintf(error,ERROR_SIZE,"%s is not a valid address!",event.data.scalar.value);
        return 1;
      }
      config->slaves[config->slaves_count-1].address = address;
    }else if(strcmp("description",event.data.scalar.value)==0){
      if(get_next_scalar_event(parser,&event,error,"Slave","Expected description")){
        return 1;
      }
      strncpy(config->slaves[config->slaves_count-1].description,event.data.scalar.value,DESCRIPTION_SIZE);
    }else if(strcmp("registries",event.data.scalar.value)==0){
      int rez = parse_registries(parser, config, error);
      if(rez!=0){
        return rez;
      }
    }
    yaml_event_delete(&event);
    if(!yaml_parser_parse(parser, &event)){
      strcpy(error,"Error parsing yaml. Slave section");
      return 1;
    }

    //printf("slave %d,%d\n",event.type,YAML_SCALAR_EVENT);

  }

  if(event.type!=YAML_MAPPING_END_EVENT){
    strcpy(error,"Error parsing yaml, no end mapping. Slave section");
    printf("%d\n",event.type );
    return 1;
  }

  return 0;
}

int parse_slaves(yaml_parser_t *parser, Configuration *config, char *error){
  yaml_event_t event;

  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Slaves section");
    return 1;
  }
  yaml_event_delete(&event);
  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Slaves section");
    return 1;
  }

  //we must wait for the final sequence event
  while(1 && event.type!=YAML_SEQUENCE_END_EVENT &&
      event.type!=YAML_DOCUMENT_END_EVENT && event.type!=YAML_NO_EVENT){
    if(event.type!=YAML_MAPPING_START_EVENT){
      strcpy(error, "Invalid yaml. Slaves section");
      return 1;
    }

    int rez = parse_slave(parser,config,error);
    if (rez!=0){
      return rez;
    }

    yaml_event_delete(&event);
    //after reading a slave we have its mapping end event.
    //we need the next event to know what to do
    if(!yaml_parser_parse(parser, &event)){
      strcpy(error,"Error parsing yaml. Slaves section");
      return 1;
    }
  }
  yaml_event_delete(&event);
  return 0;
}

//parse_connection - read the connection settings from yaml
int parse_connection(yaml_parser_t *parser, Configuration *config, char *error){
  yaml_event_t event;

  if(!yaml_parser_parse(parser, &event)){
    strcpy(error,"Error parsing yaml. Connection section");
    return 1;
  }

  //we parse through the yaml until we read all our connection
  while(1 && event.type!=YAML_STREAM_END_EVENT){
    switch(event.type){
      case YAML_SCALAR_EVENT:
        if(strcmp(event.data.scalar.value,"port")==0){//read serial port setting
          if(get_next_scalar_event(parser,&event,error,"Connection","Expected port")){
            return 1;
          }
          strncpy(config->conn.port,event.data.scalar.value,PORT_SIZE);
        }else if(strcmp(event.data.scalar.value,"speed")==0){
          if(get_next_scalar_event(parser,&event,error,"Connection","Expected speef")){
            return 1;
          }
          int speed = atoi(event.data.scalar.value);
          speed = get_speed(speed);
          if(speed ==0 ){
            snprintf(error,ERROR_SIZE,"Invalid speed detected. %s",event.data.scalar.value);
            yaml_event_delete(&event);
            return 1;
          }
          config->conn.speed = speed;
        }
        break;
      case YAML_SEQUENCE_START_EVENT:
      case YAML_MAPPING_START_EVENT:
        yaml_event_delete(&event);
        break;
      case YAML_SEQUENCE_END_EVENT:
      case YAML_MAPPING_END_EVENT:
        yaml_event_delete(&event);
        if(strlen(config->conn.port)==0 || config->conn.speed==0){
          snprintf(error,ERROR_SIZE,"Invalid value detected for port=%s and/or speed=%d",config->conn.port, config->conn.speed);
          return 1;
        }
        return 0;
        break;
      default:
        snprintf(error,ERROR_SIZE, "Unexpected yaml parse. %d",event.type);
        yaml_event_delete(&event);
        return 1;
    }
    yaml_event_delete(&event);
    if(!yaml_parser_parse(parser, &event)){
			strcpy(error,"Error parsing yaml. Connection section");
			return 1;
		}
  }
  strcpy(error,"Incomplete connection info in the configuration file.");
  return 3;
}

int parse_by_event(yaml_parser_t parser,Configuration *config, char *error){
  yaml_event_t event;
  int rez = 0;
	if(!yaml_parser_parse(&parser, &event)){
		fprintf(stderr, "%s\n", "error parsing yaml");
		return;
	}
	while(event.type!=YAML_STREAM_END_EVENT){
		yaml_event_delete(&event);
		if(!yaml_parser_parse(&parser, &event)){
			fprintf(stderr, "%s\n", "error parsing yaml");
			return;
		}

		switch (event.type) {
			case YAML_NO_EVENT:
				printf("%s\n","no event");
				break;
			case YAML_STREAM_START_EVENT:
				//printf("%s\n","STREAM START"); break;
			case YAML_STREAM_END_EVENT:
				//printf("%s\n","STREAM END"); break;
			case YAML_DOCUMENT_START_EVENT:
				//printf("%s\n","DOCUMENT START"); break;
			case YAML_DOCUMENT_END_EVENT:
			  //printf("%s\n", "DOCUMENT END"); break;
			case YAML_SEQUENCE_START_EVENT:
			  //printf("%s\n","Sequence start"); break;
			case YAML_SEQUENCE_END_EVENT:
			  //printf("%s\n","Sequence end"); break;
			case YAML_MAPPING_START_EVENT:
			  //printf("%s\n","Mapping start"); break;
			case YAML_MAPPING_END_EVENT:
			  //printf("%s\n","Mapping end"); break;
			case YAML_ALIAS_EVENT:
			  //printf("Got ALIAS: %s\n", event.data.alias.anchor ); break;
        break;
			case YAML_SCALAR_EVENT:
        if(strcmp(event.data.scalar.value,"connection")==0){
          rez = parse_connection(&parser,config,error);
        }else if(strcmp(event.data.scalar.value,"slaves")==0){
          rez = parse_slaves(&parser,config,error);
        }
        if (rez!=0){
          yaml_event_delete(&event);
          return rez;
        }
			  //printf("Got VALUE: %s\n",event.data.scalar.value); break;
			default:
       break;
			  //printf("GOT EVENT: %d\n",event.type); break;
		}
	}

	yaml_event_delete(&event);

}
