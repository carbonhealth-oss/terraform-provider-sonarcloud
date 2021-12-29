


clear:
	rm -rf .terraform.lock.hcl && terraform init

apply:
	terraform apply --auto-approve

plan:
	terraform plan

install:
	cd .. && make

reset:
	make install && make clear && make plan


