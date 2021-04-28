# Adopt default VPC into Terraform
resource "aws_default_vpc" "default" {
  tags = {
    Name        = "default-vpc"
    Environment = "Default"
  }
}

# Prevent ingress and egress from default security group
resource "aws_default_security_group" "default" {
  vpc_id = aws_default_vpc.default.id

  tags = {
    Name        = "default-sg"
    Environment = "Default"
  }
}

# Prevent ingress and egress from default network ACL
resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_default_vpc.default.default_network_acl_id

  tags = {
    Name        = "default-acl"
    Environment = "Default"
  }
}

# Remove routes from default route table
resource "aws_default_route_table" "default" {
  default_route_table_id = aws_default_vpc.default.default_route_table_id

  route = []

  tags = {
    Name        = "default-rtb"
    Environment = "Default"
  }
}

# Adopt default subnet into Terraform
resource "aws_default_subnet" "default" {
  availability_zone = "eu-west-1a"

  tags = {
    Name        = "default-subnet"
    Environment = "Default"
  }
}

# Create production VPC
resource "aws_vpc" "production" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name        = "production-vpc"
    Environment = "Production"
  }
}

# Create public subnet
resource "aws_subnet" "production_public" {
  cidr_block = "10.0.0.0/24"
  vpc_id     = aws_vpc.production.id

  tags = {
    Name        = "production-public-subnet"
    Environment = "Production"
    Access      = "Public"
  }
}

# Create private subnet
resource "aws_subnet" "production_private" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.production.id

  tags = {
    Name        = "production-private-subnet"
    Environment = "Production"
    Access      = "Private"
  }
}

# Create persistence subnet
resource "aws_subnet" "production_persistence" {
  cidr_block = "10.0.2.0/24"
  vpc_id     = aws_vpc.production.id

  tags = {
    Name        = "production-persistence-subnet"
    Environment = "Production"
    Access      = "Private"
  }
}

# Create internet gateway
# This allows access from the VPC to the internet
resource "aws_internet_gateway" "production" {
  vpc_id = aws_vpc.production.id

  tags = {
    Name        = "production-igw"
    Environment = "Production"
  }
}

# Create route table for public subnet
# Once allocated, this will allow resources in the public subnet direct access to the internet via the internet gateway
resource "aws_route_table" "production_public" {
  vpc_id = aws_vpc.production.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.production.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.production.id
  }

  tags = {
    Name        = "production-public-rtb"
    Environment = "Production"
  }
}

# Create route table association between the public subnet and the public route table
resource "aws_route_table_association" "production_public_route_table_public_subnet_association" {
  subnet_id      = aws_subnet.production_public.id
  route_table_id = aws_route_table.production_public.id
}

# Create elastic IP address
# Once associated to the NAT Gateway, this allows communication from the internet to the VPC
resource "aws_eip" "production" {
  depends_on = [
    aws_internet_gateway.production
  ]

  vpc = true

  tags = {
    Name        = "production-elastic-ip"
    Environment = "Production"
  }
}

# Create NAT gateway
# This will be used to allow resources from private subnets to communicate with the internet
# It lives in the public subnet
resource "aws_nat_gateway" "production" {
  allocation_id = aws_eip.production.id
  subnet_id     = aws_subnet.production_public.id

  tags = {
    Name        = "production-nat-gateway"
    Environment = "Production"
  }
}

# Create route table for private subnet
# This allows resources in the private subnet to communicate with the internet via the NAT Gateway
resource "aws_route_table" "production_private" {
  vpc_id = aws_vpc.production.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.production.id
  }

  tags = {
    Name        = "production-private-rtb"
    Environment = "Production"
  }
}

# Create route table association between the private route table and the private subnet
resource "aws_route_table_association" "production_private_route_table_private_subnet_association" {
  subnet_id      = aws_subnet.production_private.id
  route_table_id = aws_route_table.production_private.id
}

# Create route table association between the public subnet and the public route table
resource "aws_route_table_association" "production_private_subnet_internet_access" {
  subnet_id      = aws_subnet.production_public.id
  route_table_id = aws_route_table.production_public.id
}

resource "aws_security_group" "production_lambda_sg" {
  vpc_id = aws_vpc.production.id
  name   = "production-lambda-sg"

  ingress {
    from_port = 443
    protocol  = "tcp"
    to_port   = 443
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    ipv6_cidr_blocks = [
      "::/0"
    ]
  }

  egress {
    from_port = 0
    protocol  = "tcp"
    to_port   = 65535
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    ipv6_cidr_blocks = [
      "::/0"
    ]
  }

  tags = {
    Name        = "production-lambda-sg"
    Environment = "Production"
  }
}

resource "aws_security_group" "production_elasticache_sg" {
  vpc_id = aws_vpc.production.id
  name   = "production-elasticache-sg"

  ingress {
    description = "ElastiCache Redis access"
    from_port   = 6379
    protocol    = "tcp"
    to_port     = 6379
    security_groups = [
      aws_security_group.production_lambda_sg.id
    ]
  }

  tags = {
    Name        = "production-elasticache-sg"
    Environment = "Production"
  }
}
