package com.cardinfolink.yunshouyin.data;

public class User {
 private String username;
 private String password;
 private String bank_open;
 private String payee;
 private String payee_card;
 private String phone_num;
 private String clientid;
 private String email;
 private String object_id;
 
 
 
 
 
 public String getObject_id() {
	return object_id;
}


public void setObject_id(String object_id) {
	this.object_id = object_id;
}


public String getEmail() {
	return email;
}


public void setEmail(String email) {
	this.email = email;
}
private String limit="true";
 
public String getLimit() {
	return limit;
}


public void setLimit(String limit) {
	this.limit = limit;
}


public String getClientid() {
	return clientid;
}


public void setClientid(String clientid) {
	this.clientid = clientid;
}


public String getUsername() {
	return username;
}
public void setUsername(String username) {
	this.username = username;
}
public String getPassword() {
	return password;
}
public void setPassword(String password) {
	this.password = password;
}
public String getBank_open() {
	return bank_open;
}
public void setBank_open(String bank_open) {
	this.bank_open = bank_open;
}
public String getPayee() {
	return payee;
}
public void setPayee(String payee) {
	this.payee = payee;
}
public String getPayee_card() {
	return payee_card;
}
public void setPayee_card(String payee_card) {
	this.payee_card = payee_card;
}
public String getPhone_num() {
	return phone_num;
}
public void setPhone_num(String phone_num) {
	this.phone_num = phone_num;
}
 

}
