


type Identifier struct {
  name string
  id string
}

type Area struct {
  Identifier
}

type VirtualMachine struct {
  Identifier
}

type DnsService struct {
  Identifier
  VirtualMachine
}

type CloudBase struct {
  Identifier
  DnsService
} 
