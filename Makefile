NAME		=	ib
GATEWAY_IMAGE	=	$(NAME)_gateway_test
GATEWAY_CONT	=	$(GATEWAY_IMAGE)_c
TEST_CONT	=	$(NAME)_test_c

all		:	test

.build_gw_id	:	testserver
		docker build -t $(GATEWAY_IMAGE) testserver
		@docker inspect -f '{{.Id}}' $(GATEWAY_IMAGE) > .build_gw_id

.gateway_id	:	.build_gw_id
		-@docker rm -f $(GATEWAY_CONT) > /dev/null 2> /dev/null || true
		-@docker rm -f $(GATEWAY_CONT)_tmp > /dev/null 2> /dev/null || true
		docker run --name $(GATEWAY_CONT) -d $(GATEWAY_IMAGE)
		@echo Wait for Gateway to be started
		@sleep 1
		@docker run --link $(GATEWAY_CONT):gw --rm --name $(GATEWAY_CONT)_tmp -t ubuntu:14.04 \
			bash -c 'for i in {1..60}; do \
					echo | nc $$GW_PORT_4003_TCP_ADDR 4002 && exit 0 || (echo -n ..; sleep 1); \
				done; \
				echo; \
				echo Waiting for Gateway timed out; exit 1'
		@echo
		@docker inspect -f '{{.Id}}' $(GATEWAY_IMAGE) > .gateway_id

gateway		:	.gateway_id

.build_id	:	.
		docker build -t $(NAME) .
		@docker inspect -f '{{.Id}}' $(GATEWAY_IMAGE) > .build_id

build		:	.build_id

test		:	gateway build
		-@docker rm -f $(TEST_CONT) > /dev/null 2> /dev/null || true
		docker run --link $(GATEWAY_CONT):gw --name $(TEST_CONT) -t $(NAME) bash -c 'cd /src && go test $(TESTFLAGS) -gw $$GW_PORT_4003_TCP_ADDR:4003'

clean		:
		-@docker rm -f $(GATEWAY_CONT) > /dev/null 2> /dev/null || true
		-@docker rm -f $(GATEWAY_CONT)_tmp > /dev/null 2> /dev/null || true
		-@docker rm -f $(TEST_CONT) > /dev/null 2> /dev/null || true

clean_all	:	clean
		-@rm -f .build_id .build_gw_id .gateway_id

re		:	clean_all all


.PHONY		:	all gateway buld test clean clean_all re
