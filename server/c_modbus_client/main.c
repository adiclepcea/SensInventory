#include <stdio.h>
#include <stdlib.h>
#include <termios.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include "configuration_parser.h"

int setup_interface(char *portname, int speed, int parity, char *error ){

  int fd = open(portname, O_RDWR | O_NOCTTY | O_NDELAY);

  if (fd ==-1){ //port open failed
    snprintf(error,ERROR_SIZE, "%s",strerror(errno)); //we store the error
    return fd;
  }

  struct termios tty;
  memset(&tty,0,sizeof(tty));

  if(tcgetattr(fd,&tty)!=0){
    snprintf(error,ERROR_SIZE, "%s from tcgetattr",strerror(errno)); //we store the error
    return -1;
  }

  cfsetospeed(&tty,speed);
  cfsetispeed(&tty,speed);

  tty.c_cflag = (tty.c_cflag & ~CSIZE) | CS8;
  tty.c_iflag &= ~IGNBRK; //no break processing
  tty.c_lflag = 0; //no echo, signal chars etc.
  tty.c_oflag = 0; //no delays, remappings
  tty.c_cc[VMIN]  = 0; //read doesn't block
  tty.c_cc[VTIME] = 5; //0.5 seconds read timeout
  tty.c_iflag &= ~(IXON | IXOFF | IXANY); // shut off xon/xoff ctrl
  tty.c_cflag |= (CLOCAL | CREAD);// ignore modem controls, enable reading
  tty.c_cflag &= ~(PARENB | PARODD);      // shut off parity
  tty.c_cflag |= parity;
  tty.c_cflag &= ~CSTOPB;
  tty.c_cflag &= ~CRTSCTS;

  if (tcsetattr (fd, TCSANOW, &tty) != 0)
  {
    snprintf(error,ERROR_SIZE, "%s from tcsetattr",strerror(errno)); //we store the error
    return -1;
  }

  if (set_blocking(fd,0,error)!=0){
    close(fd);
    return -1;
  }

  return fd;
}

int set_blocking (int fd, int should_block, char *error)
{
  struct termios tty;
  memset (&tty, 0, sizeof tty);
  if (tcgetattr (fd, &tty) != 0)
  {
    snprintf (error,ERROR_SIZE, "error %d from tggetattr", errno);
    return -1;
  }

  tty.c_cc[VMIN]  = should_block ? 1 : 0;
  tty.c_cc[VTIME] = 5;            // 0.5 seconds read timeout

  if (tcsetattr (fd, TCSANOW, &tty) != 0){
    snprintf (error,ERROR_SIZE,"error %d setting term attributes", errno);
    return -1;
  }
  return 0;
}

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

  //print_config(&config);

  int fd = setup_interface(config.conn.port,config.conn.speed,0,error); //change second param for speed
  if (fd==-1){ //error
    fprintf(stderr, "failed to setup and open interface %s: %s\n",config.conn.port, error);
    free(error);
    return 1;
  }else{
    write(fd,"hello!\r\n",strlen("hello!\r\n"));
    usleep ((7 + 25) * 100);
    close(fd);
  }

  return 0;
}
