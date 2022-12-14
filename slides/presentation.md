class: center, middle

<div class="my-header"><img src="assets/imc-logo.png" style="float: left; width:250px"></div>

.image-40[![rimworld](assets/rimworld-dc.png)]
# The 16-bit Datacenter
Brandon Ewing<br />
CHI-NOG 10, Oct 2022<br />
[www.github.com/bewing/chinog-10](https://www.github.com/bewing/chinog-10/)

???
* Thank the PC
* intro self
* reducing complexitiy in datacenter configuration
* complexity means more data collection and entry
* more touches means more opportunity for error
* Examples today in Aristese due to my familiarity

---
<div class="my-header"><h1>Talk Contents</h1></div>
<br />
.row.table.middle[
.col-6[
<div style="
        border-radius: 100%;
        width: 400px;
        aspect-ratio: 1;
        background: conic-gradient(red 0deg 144deg, green 144deg 288deg, orange 240deg 360deg
            );
    "></div>
]
.col-6[
.big.red[40% Knowledge Transfer]<p>
.big.green[40% Stuff to think about]<br />
]]

--

.row.table.middle[
.col-6[&nbsp;]
.col-6[
.big.orange[20% Please tell me he's not doing that in prod]<br />
]
]

---
class: middle
<div class="my-header"><h1>Simplicity</h1></div>
.image-40[![one-dc](assets/one-dc.png)]

???
Simple, Spine/Leaf datacenter with service leaves

---
class: full, middle, center
background-image: url(assets/months_later.png)

---
class: middle

<div class="my-header"><h1>Complexity</h1></div>
.image-60[![horizontal](assets/complex.png)]

???
Datacenters spread across regions
different deploy times
different hardware
even different vendors!
---

<div class="my-header"><h1>Complexity</h1></div>

.image-60[![one-arrow](assets/one-arrow.png)]

---

<div class="my-header"><h1>Complexity</h1></div>

.image-60[![two-arrow](assets/two-arrow.png)]

--

```terminal
$ wc -l configs/DC1-LEAF1
2130 configs/DC1-LEAF1
$ git diff --color --stat --no-index configs/DC1-LEAF1 configs/DC4-LEAF9
 configs/{DC1-LEAF1 => DC4-LEAF9} | 881 <span style="color:lime;"> +++++++++ </span><span style="color:red;">-----------</span>
 
 1 file changed, 362 insertions(+), 519 deletions(-)

```

---
class: middle

<div class="my-header"><h1>Baseline</h1></div>

* Use templates
* Same templates for provisioning and configuration
* If you aren't pushing configs, at least run audits
* Clean up cruft!

???
Templating language doesn't matter<br />
I'm using gomplate from containerlab here

---
class: inverse
<div class="my-header"><h1>Context</h1></div>

```yaml
nodes:
- hostname: dc1-leaf1
  region: dc1
  role: leaf
  aaa:
  - 192.168.3.2
  - 192.168.10.1
  - 192.168.14.10
  addresses:
    Loopback0:
      ipv4:
      - 172.18.4.10/32
      ipv6:
      - "2001:db8::1/128"
    Ethernet1/1:
      ipv4:
      - 10.0.0.0/31
      ipv6:
      - "2001:2b8:1::1/64"
  bgp:
    router-id: 172.18.4.10
    asn: "65534"
    ipv4:
      peer-groups:
      - name: SPINE-LEAF
        outbound-policy: SPINE-TO-LEAF
        inbound-policy: LEAF-TO-SPINE
        neighbors:
        - address: 10.0.0.1
          asn: "65234"
  dns:
    search: warningg.com
    servers:
    - 8.8.8.8
```

---
class:middle, center
.big[
```terminal
$ stat --format "%s" host.yaml
608
```
]
--
class: inverse

.big[
608 bytes > 16 bits
]

---
class: middle
<div class="my-header"><h1>Boilerplate</h1></div>

.table[
.row[
.col-6[
.big[DC1-LEAF4]
]
.col-6[
.big[DC4-LEAF9]
]
]
.row[
.col-6[
```terminal
<span style="color:red;">-ip nameserver 192.168.0.0</span>
<span style="color:red;">-ip nameserver 192.168.100.100</span>
!
<span style="color:red;">-tacacs-server 192.168.3.2</span>
tacacs-server 192.168.10.1
<span style="color:red;">-tacacs-server 192.168.14.10</span>

```
]
.col-6[
```terminal
<span style="color:green;">+ip nameserver 10.240.0.0</span>
<span style="color:green;">+ip nameserver 10.244.100.100</span>
!
<span style="color:green;">+tacacs-server 192.168.14.10</span>
tacacs-server 192.168.10.1
<span style="color:green;">+tacacs-server 192.168.3.2</span>
```
]
]
]

---
class: middle
<div class="my-header"><h1>Boilerplate</h1></div>

.image-60[![anycast](assets/anycast.png)]

---
class: middle
<div class="my-header"><h1>Boilerplate</h1></div>

.table[
.row[
.col-6[
.big[DC1-LEAF4]
]
.col-6[
.big[DC4-LEAF9]
]
]
.row[
.col-6[
```terminal
ip nameserver 192.168.100.100
!
tacacs-server 192.168.3.2
tacacs-server 192.168.10.1
tacacs-server 192.168.14.10

```
]
.col-6[
```terminal
ip nameserver 192.168.100.100
!
tacacs-server 192.168.3.2
tacacs-server 192.168.10.1
tacacs-server 192.168.14.10
```
]
]
]
* Anycast when possible
 * Make sure your application withdraws itself if not healthy
 * Avoid ECMP through policy
* Global services with fallbacks if not
 * If RTT is important, programmatically order them with template logic

???
You don't have to golf if you don't want to

---
class: inverse
<div class="my-header"><h1>Context</h1></div>
<br />
```yaml
nodes:
- hostname: dc1-leaf1
  region: dc1
  role: leaf
  addresses:
    Loopback0:
      ipv4:
      - 172.18.4.10/32
      ipv6:
      - "2001:db8::1/128"
    Ethernet1/1:
      ipv4:
      - 10.0.0.0/31
      ipv6:
      - "2001:2b8:1::1/64"
  bgp:
    router-id: 172.18.4.10
    asn: "65534"
    ipv4:
      peer-groups:
      - name: SPINE-LEAF
        outbound-policy: SPINE-TO-LEAF
        inbound-policy: LEAF-TO-SPINE
        neighbors:
        - address: 10.0.0.1
          asn: "65234"
```

---
# Don't Repeat Yourself
* Don't have the same info twice (IPv4 Lo0, router-id)
* Generate the IPv6 loopback

<br />
.biggish[
2001:db8::/64 + 172.18.4.10/32<br /><br />
.col-6[2001:db8::ac12:40ac]
.col-6[2001:db8::172.18.4.10]
.col-4[2001:db8::172:18:4:10]
]

---
class: inverse
<div class="my-header"><h1>Context</h1></div>
<br />
```yaml
nodes:
- hostname: dc1-leaf1
  region: dc1
  role: leaf
  router-id: 172.18.4.10
  addresses:
    Ethernet1/1:
      ipv4:
      - 10.0.0.0/31
      ipv6:
      - "2001:2b8:1::1/64"
  bgp:
    asn: "65534"
    ipv4:
      peer-groups:
      - name: SPINE-LEAF
        outbound-policy: SPINE-TO-LEAF
        inbound-policy: LEAF-TO-SPINE
        neighbors:
        - address: 10.0.0.1
          asn: "65234"
```

???
Focus on region / role / routerid <br />
Ask who's used BGP info communities before
What is a router-id?   32 bits<br />

---
<div class="right-header" text-align="right">This is the 20% part</div>

# Don't Repeat Yourself
* Don't have the same info twice (IPv4 Lo0, router-id)
* Generate the IPv6 loopback

.row.table[
.col-3[
.center[172]
]
.col-3[
.center[18]
]
.col-3[
.center[4]
]
.col-3[
.center[10]
]]
.row.table.middle[
.col-3[
.center[10101100]
]
.col-3[
.center[00010010]
]
.col-3[
.center[00000100]
]
.col-3[
.center[00001010]
]]

???
use the lower 16 bits to store role info<br />

--

.row.table.middle[
.col-3[]
.col-5[
4
]
.col-3[
10
]]
.row.table.middle[
.col-2[
.orange[00]
]
.col-2[
.green[00]
]
.col-2[
.red[01]
]
.col-2[
.purple[00]
]
.col-6[
.purple[00001010]
]]

.row.table.middle[
.col-2[
.orange[Region]
]
.col-2[
.green[Site]
]
.col-2[
.red[Layer]
]
.col-3[
.purple.center[Device]
]
]

---
<div class="right-header" text-align="right">This is the 20% part</div>

# Don't Repeat Yourself
* Don't have the same info twice (IPv4 Lo0, router-id)
* Generate the IPv6 loopback
* Encode information into bit fields
* BGP ASN?  Sure!
<br />
<br />

.row.table.middle[
.col-2[
.orange[Region]
]
.col-2[
.green[Site]
]
.col-2[
.red[Layer]
]
.col-3[
.purple.center[Device]
]
]

.center[.big[
65
.orange[X]
.green[X]
.red[X]
.
.orange[X]
.green[X]
.red[X]
.purple[XX]
]]
.center[.big[**<16 bits>.<16 bits>**]]

---
class: middle
<div class="right-header" text-align="right">This is the 20% part</div>

.col-6[
.center.big[
REGION 0 LEAF
]
]
.col-6[
.center.big[
172.18.4.10
]
]

.col-2[
.orange[Region]
]
.col-2[
.green[Site]
]
.col-2[
.red[Layer]
]
.col-2[&nbsp;]
.col-4[
.purple[Device]
]

.col-2[
.orange[00]
]
.col-2[
.green[00]
]
.col-2[
.red[01]
]
.col-2[
.purple[00]
]
.col-4[
.purple[00001010]
]
.col-6[
.black[65]
.orange[0]
.green[0]
.red[1] .
.black[(4 * 256 + 10)]
]
.col-6[
**65001.1034**
]

---
class: middle
<div class="right-header" text-align="right">This is the 20% part</div>

.col-6[
.center.big[
REGION 3 SUPERSPINE
]
]
.col-6[
.center.big[
172.18.72.75
]
]

.col-2[
.orange[Region]
]
.col-2[
.green[Site]
]
.col-2[
.red[Layer]
]
.col-2[&nbsp;]
.col-4[
.purple[Device]
]

.col-2[
.orange[11]
]
.col-2[
.green[00]
]
.col-2[
.red[11]
]
.col-2[
.purple[00]
]
.col-4[
.purple[01001011]
]
.col-6[
.black[65]
.orange[3]
.green[0]
.red[3]
.black[.(72 * 256 + 75)]
]
.col-6[
**65303.18507**
]

---
class: middle
<div class="right-header" text-align="right">This is the 20% part</div>

.row.table.middle[
.col-6[
* Why 16 bits?
 * Have to title the talk somehow
 * IPv4 Reachability
 * BGP 4-byte Private ASN space ~25 bits
]
.col-6[
.orange[Please do not do this in prod]
]
]

---
class: inverse
<div class="my-header"><h1>Context</h1></div>
<br />
```yaml
nodes:
- router-id: 172.18.4.10
  addresses:
    Ethernet1/1:
      ipv4:
      - 10.0.0.0/31
      ipv6:
      - "2001:2b8:1::1/64"
  bgp:
    ipv4:
      peer-groups:
      - name: SPINE-LEAF
        outbound-policy: SPINE-TO-LEAF
        inbound-policy: LEAF-TO-SPINE
        neighbors:
        - address: 10.0.0.1
          asn: "65234"
```

---
class: middle
<div class="my-header"><h1>Interfaces</h1></div>

.table[
.row[
.col-6[
DC1-LEAF4
]
.col-6[
DC4-LEAF9
]
]
.row[
.col-6[
```terminal
interface Ethernet49/1
<span style="color:red;">-  description DC1-SPINE1:Et5/1</span>
  no switchport
<span style="color:red;">-  ip address 10.0.0.0/31</span>
<span style="color:red;">-  ipv6 address 2001:db8::ffff:0a00:0/127</span>
  pim ipv4 sparse-mode
  pim ipv6 sparse-mode
```
]
.col-6[
```terminal
interface Ethernet 49/1
<span style="color:green;">+  description DC4-SPINE3:Et5/1</span>
  no switchport
<span style="color:green;">+  ip address 10.5.49.22/31</span>
<span style="color:green;">+  ipv6 address 2001:db8::ffff:a05:3116/127</span>
  pim ipv4 sparse-mode
  pim ipv6 sparse-mode
```
]
]
]

.row.table.middle[
.col-6[
```terminal
router bgp 65001.1034
  neighbor SPINES-v4 peer-group
  neighbor SPINES-v6 peer-group
<span style="color:red;">-  neighbor 10.0.0.1 peer-group SPINES-v4</span>
<span style="color:red;">-  neighbor 10.0.0.1 remote-as 65002.3080</span>
<span style="color:red;">-  neighbor 2001:db8:ffff:0a00:1 peer-group SPINES-v6</span>
<span style="color:red;">-  neighbor 2001:db8:ffff:0a00:1 remote-as 65002.3080</span>
  address-family ipv4 unicast
    peer-group SPINES-v4 activate
    no peer-group SPINES-v6 activate
  !
  address-family ipv6 unicast
    no peer-group SPINES-v4 activate
    peer-group SPINES-v6 activate
  !
!
```
]
.col-6[
```terminal
router bgp 65031.1077
  neighbor SPINES-v4 peer-group
  neighbor SPINES-v6 peer-group
<span style="color:green;">+  neighbor 10.5.49.23 peer-group SPINES-v4</span>
<span style="color:green;">+  neighbor 10.5.49.23 remote-as 65032.3088</span>
<span style="color:green;">+  neighbor 2001:db8::ffff:a05:3117 peer-group SPINES-v6</span>
<span style="color:green;">+  neighbor 2001:db8::ffff:a05:3117 remote-as 65032.3088</span>
  address-family ipv4 unicast
    peer-group SPINES-v4 activate
    no peer-group SPINES-v6 activate
  !
  address-family ipv6 unicast
    no peer-group SPINES-v4 activate
    peer-group SPINES-v6 activate
  !
!
```
]
]

???
Descriptions:  operator-dependent.  Some people love them

---
class: middle
<div class="my-header"><h1>Interfaces</h1></div>

* [RFC5549](https://datatracker.ietf.org/doc/html/rfc5549) - IPv4 NLRI in IPv6 Peering
* Allows IPv4 reachability over just IPv6 peerings
* No longer need IPv4 BGP peerings
* No longer need IPv4 point to point interfaces!

```terminal
interface Ethernet49/1
  ipv6 address 2001:db8::ffff:0a00:0/127
  pim ipv4 sparse-mode
  pim ipv6 sparse-mode
!
router bgp 65001.1034
  neighbor SPINES peer-group
  neighbor 2001:db8::ffff:0a00:1/127 peer-group SPINES
  neighbor 2001:db8::ffff:0a00:1/127 remote-as 65002.3080
  address-family ipv4 unicast
    bgp next-hop address-family ipv6
    neighbor SPINES activate
    neighbor SPINES next-hop address-family ipv6 originate
  !
  address-family ipv6 unicast
    neighbor SPINES activate
  !
!
```

---
class: middle
<div class="my-header"><h1>But wait, there's more!</h1></div>

* [draft-white-linklocal-capability](https://datatracker.ietf.org/doc/draft-white-linklocal-capability/)
* Widely supported across vendors
* Standardizes existing practice of BGP peering via link-local IPv6 addresses
* Now we don't need any globally unique addressing!

```terminal
interface Ethernet49/1
  ipv6 address fe80::0/64
  pim ipv4 sparse-mode
  pim ipv6 sparse-mode
!
router bgp 65001.1034
  neighbor SPINES peer-group
  neighbor fe80::1%Ethernet49/1 peer-group SPINES
  neighbor fe80::1%Ethernet49/1 remote-as 65002.3080
  address-family ipv4 unicast
    bgp next-hop address-family ipv6
    neighbor SPINES activate
    neighbor SPINES next-hop address-family ipv6 originate
  !
  address-family ipv6 unicast
    neighbor SPINES activate
  !
!
```

---
class: middle

<div class="my-header"><h1>BGP Peerings</h1></div>
* BGP Peer autodetection
* Multiple IETF IDR WG proposals (see [draft-ietf-idr-bgp-autoconf-considerations](https://datatracker.ietf.org/doc/draft-ietf-idr-bgp-autoconf-considerations/02/))
* Some Layer2 (LLDP), some Layer3
* Some are secured, some aren't
* Some are stateful, some are stateless
* No real consensus yet
* Trying to support all DC use cases

---
class: middle
<div class="my-header"><h1>BGP Peerings</h1></div>

* More than one vendor supports IPv6 link-local peer autodetection
* Uses IPv6 RAs to identify routers on an interface
* No current RFC or I-D for this behavior
* There may be interoperability issues

```terminal
interface Ethernet49/1
  ipv6 enable
  pim ipv4 sparse-mode
  pim ipv6 sparse-mode
!
peer-filter SPINES
 10 match 4200000000-4294967294 result accept
!
router bgp 65001.1034
  neighbor SPINES peer-group
  neighbor interface Ethernet49/1 peer-group SPINES peer-filter SPINES
  address-family ipv4 unicast
    bgp next-hop address-family ipv6
    neighbor SPINES activate
    neighbor SPINES next-hop address-family ipv6 originate
  !
  address-family ipv6 unicast
    neighbor SPINES activate
  !
!
```

---
class: middle
<div class="my-header"><h1>BGP Peerings</h1></div>
* Neighborships can specify AS ranges
* Still have to statically map a policy to an interface

** DC1-SPINE1 **
```terminal
neighbor interface Ethernet1-48 peer-group LEAVES peer-filter LEAVES
neighbor interface Ethernet49/1-52/1 peer-group SUPERSPINES peer-filter SUPERSPINES
```
* While we can predict this in our lab, production isn't as nice

```terminal
neighbor interface Ethernet52/4 peer-group LEAVES peer-filter TEMP-LEAF
```

---
class:middle
<div class="my-header"><h1>Hosts</h1></div>
* If your virtualization and tenancy model supports it, go deeper!
* Assign hosts router-ids
 * Deployment
 * Custom DHCP
* Advertise reachability via host-based BGP
 * GoBGP
 * FRR
* IPv6 LL Autodetection Supported!
 
---
<div class="my-header"><h1>Hosts</h1></div>

.image-40[![one-dc](assets/hosts.png)]

* Connect servers to one or more leaves
* No more Layer 2 problems!
 * LACP/MLAG
 * Spanning Tree
---

<div class="my-header"><h1>Hosts</h1></div>
<br />

```terminal
router bgp 65001.1034
  neighbor SPINES peer-group
  neighbor SERVERS peer-group
  neighbor interface Ethernet1-48 peer-group SERVERS peer-filter SERVERS
  neighbor interface Ethernet49/1-52/4 peer-group SPINES peer-filter SPINES
  address-family ipv4 unicast
    bgp next-hop address-family ipv6
    neighbor SERVERS activate
    neighbor SERVERS next-hop address-family ipv6 originate
    neighbor SPINES activate
    neighbor SPINES next-hop address-family ipv6 originate
  !
  address-family ipv6 unicast
    neighbor SERVERS activate
    neighbor SPINES activate
  !
!
```

--

* What would be great is if remote AS determined policy
* Ask your vendor if this is something they can support
* May require stronger security ([RFC 5925 TCP-AO](https://datatracker.ietf.org/doc/html/rfc5925))

---
# Recap
* Eliminated boilerplate
 * templates
 * anycast
* Derived as much as possible from the router-id
* Removed non-loopback addressing
* BGP to everything

---
class: middle, inverse
<div class="my-header"><h1>Context</h1></div>

.row.table.middle[
.col-6[
```yaml
nodes:
- hostname: dc1-leaf1
  region: dc1
  role: leaf
  aaa:
  - 192.168.3.2
  - 192.168.10.1
  - 192.168.14.10
  addresses:
    Loopback0:
      ipv4:
      - 172.18.4.10/32
      ipv6:
      - "2001:db8::1/128"
    Ethernet1/1:
      ipv4:
      - 10.0.0.0/31
      ipv6:
      - "2001:2b8:1::1/64"
  bgp:
    router-id: 172.18.4.10
    asn: "65534"
    ipv4:
      peer-groups:
      - name: SPINE-LEAF
        outbound-policy: SPINE-TO-LEAF
        inbound-policy: LEAF-TO-SPINE
        neighbors:
        - address: 10.0.0.1
          asn: "65234"
  dns:
    search: warningg.com
    servers:
    - 8.8.8.8
```
]
.col-6[
```yaml
nodes:
- router-id: 172.18.8.0     # site1-spine
- router-id: 172.18.8.1     # site1-spine
- router-id: 172.18.4.2     # site1-leaf
- router-id: 172.18.4.3     # site1-leaf
- router-id: 172.18.4.4     # site1-leaf
- router-id: 172.18.0.5     # site1-server
- router-id: 172.18.0.6     # site1-server
- router-id: 172.18.0.7     # site1-server
- router-id: 172.18.0.8     # site1-server
- router-id: 172.18.28.9    # site2-superspine
- router-id: 172.18.24.10   # site2-spine
- router-id: 172.18.20.11   # site2-leaf
- router-id: 172.18.16.12   # site2-server
```

]]

---
<div id="generateDemo"></div>

---
<div id="containerlabDemo"></div>


---
class: center, middle
<div class="my-header"><img src="assets/imc-logo.png" style="float: left; width:250px"></div>

# Thank you <br />
## Questions ? <br /><br />
[www.github.com/bewing/chinog-10](https://www.github.com/bewing/chinog-10/)

---

# Bitshifting
```golang
func loadNodeData(routerId string) (NodeData, error) {
	nd := NodeData{}
	ip, err := netip.ParseAddr(routerId)
	if err != nil {
		return nd, err
	}
	data := ip.AsSlice()[2]
	typeByte := data & 12 >> 2
	if typeByte^1 == 0 {
		nd.Type = "leaf"
		nd.Layer = 1
	} else if typeByte^2 == 0 {
		nd.Type = "spine"
		nd.Layer = 2
	} else if typeByte^3 == 0 {
		nd.Type = "superspine"
		nd.Layer = 3
	} else {
		nd.Type = "server"
		nd.Layer = 0
	}
	nd.Region = int(data & 192 >> 6)
	nd.Site = int(data & 48 >> 4)

	nd.ASN = fmt.Sprintf("65%d%d%d.%d", nd.Region, nd.Site, nd.Layer, int(ip.AsSlice()[2])*256+int(ip.AsSlice()[3]))
	return nd, nil
}
```
