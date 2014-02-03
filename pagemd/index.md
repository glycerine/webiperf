
<div class="container">
<div class="span-48">

<h1>iPerf Traffic Generator Front-End</h1>

<form action="../configure-iperf">

<input type="radio" name="trafprot" value="tcp" checked="checked">TCP</input>
<input type="radio" name="trafprot" value="udp">UDP</input>
<hr>

<input type="radio" name="ipv" value="ipv4" checked="checked">IPv4</input>
<input type="radio" name="ipv" value="ipv6">IPv6</input>
<hr>
<h2>Server</h2>
</p>
IPv4 address: 
<input name="Sv4o1" type="number" min="1" max="255" maxlength="3" placeholder="192" onkeypress="return isNumberKey(event)" class="ipbyte"></input><b>.</b>
<input name="Sv4o2" type="number" min="1" max="255" maxlength="3" placeholder="168" onkeypress="return isNumberKey(event)" class="ipbyte"></input><b>.</b>
<input name="Sv4o3" type="number" min="1" max="255" maxlength="3" placeholder="0" onkeypress="return isNumberKey(event)" class="ipbyte"></input><b>.</b>
<input name="Sv4o4" type="number" min="1" max="255" maxlength="3" placeholder="10" onkeypress="return isNumberKey(event)" class="ipbyte"></input>

<br>
IPv6 address: 
<input type="text" maxlength="2" class="ip6byte" name="Sv6o1"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o2"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o3"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o4"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o5"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o6"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o7"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Sv6o8"></input> /
<input type="number" maxlength="3" class="ip6byte" name="Sv6o9"></input>
<hr>

<h2>Client</h2>
</p>
IPv4 address: 
<input type="number" min="1" max="255" width="24" name="Cv4o1"></input>.
<input type="number" min="0" max="255" width="24" name="Cv4o2"></input>.
<input type="number" min="0" max="255" width="24" name="Cv4o3"></input>.
<input type="number" min="0" max="255" width="24" name="Cv4o4"></input>
<hr>
IPv6 address: 
IPv6 address: 
<input type="text" maxlength="2" class="ip6byte" name="Cv6o1"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o2"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o3"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o4"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o5"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o6"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o7"></input> :
<input type="text" maxlength="2" class="ip6byte" name="Cv6o8"></input> /
<input type="number" maxlength="3" class="ip6byte" name="Cv6o9"></input>
<hr>

<input type="checkbox" name="trafprot" value="tcp">TCP</input>
<input type="checkbox" name="trafprot" value="udp">UDP</input>
<hr>



select-format:
<select name="select-format">
   <option value="kbit">kbit</option>
   <option value="Mbit">Mbit</option>
</select>

<hr>
input data file: <input type="file" name="data-input-file"></input>

</form>

</div>
</div>
