hostname {{ .RouterId }}
{{ $nodeType := nodeType .RouterId }}
{{ $asn := asnFromRouterId .RouterId }}
{{ $nodeType }}
foo
username admin privilege 15 secret admin
!
service routing protocols model multi-agent
!
ip routing
ipv6 unicast-routing
ip route vrf MGMT 0.0.0.0/0 192.0.2.254
ipv6 route vrf MGMT ::/0 2001:fb8::1
!
vrf MGMT
!
interface Management1
   vrf MGMT
   ip address 192.0.2.0/24
   ipv4 address 2001:fb8::10/64
!
interface Loopback0
   ip address {{ .RouterId }}/32
!
interface range Ethernet1-52/1
   no switchport
   ipv6 enable
   pim ipv4 sparse-mode
   pim ipv6 sparse-mode
!
management api gnmi
   transport grpc default
   vrf MGMT
!
management api netconf
   transport ssh default
   vrf MGMT
!
management api http-commands
   no shutdown
   vrf MGMT
     no shutdown
   !
!
router bgp 65500.{{ $asn }}
   bgp asn notation asdot
   router-id {{ .RouterId }}
   neighbor SPINES peer-group
   neighbor LEAVES peer-group
   neighbor SERVERS peer-group
   {{- if eq $nodeType "leaf" }}
   neighbor interface Et1-48 peer-group SERVERS peer-filter SERVERS
   neighbor interface Et49/1-52/1 peer-group SPINES peer-filter SPINES
   {{- end }}
   !
   address-family ipv4
      bgp next-hop address-family ipv6
      neighbor SPINES activate
      neighbor SPINES next-hop address-family ipv6 originate
      neighbor LEAVES activate
      neighbor LEAVES next-hop address-family ipv6 originate
      neighbor SERVERS activate
      neighbor SERVERS next-hop address-family ipv6 originate
      redistribute connected
   !
   address-family ipv6
      neighbor SPINES activate
      neighbor LEAVES activate
      neighbor SERVERS activate
      redistribute connected
   !
end
