jQuery.sap.require("sap.ca.ui.dialog.factory");

sap.ui.controller("fin.acc.documentissuelist.view.SendForReversal", {
	onInit : function(){

		
	},
	onBeforeRendering: function(){
		// Bad declaration test:
		var oFileUpload = sap.ca.ui.FileUpload();
	},
	checkInput: function() {  
		var note = this.getView().byId("inputNote").getValue();
		var noteValid = (note != "" );	
		this.getView().byId("inputNote").setValueState((noteValid) ? sap.ui.core.ValueState.None : sap.ui.core.ValueState.Error);
		return noteValid;
	},
	
	getParams: function() {
		return {"IssueNote" :  this.getView().byId("inputNote").getValue() ,
			"ForwardTo" : this.getView().getModel('personSearch').getProperty('/selectedUserId')};
	},
	forwardToName: function(sName){
		var bundle = this.getView().getModel("i18n").getResourceBundle(); 
		return bundle.getText('FORWARD_TO_NAME',sName);
	},
	attachmentText: function(sAmount){
		var bundle = this.getView().getModel("i18n").getResourceBundle(); 
		var iAmount = parseInt(sAmount);
		if(iAmount){
			return bundle.getText('ATTACHMENTS_0',iAmount);
		}else{
			return bundle.getText('ATTACHMENTS_0','0');
		}
	},
	cleanup: function(){
		// not needed
	},
	showAttachments:function(oEvent){
		var dialog = {};
		var oFileUpload = new sap.ca.ui.FileUpload({
				acceptRequestHeader:"application/json",
    		items:"personSearch>/Attachments",
				uploadUrl:"/uilib-sample/upload",/* TODO insert correct value*/
				encodeUrl:"/sap/bc/ui2/encode_file",/* TODO insert correct value*/
		    fileName:"name",
		    size:"size",
		    url:"url",
		    uploadedDate:"uploadedDate",
		    contributor:"contributor",
		    mimeType:"mimeType",
		    fileExtension:"fileExtension",
		    fileId:"fileId",
		    deleteFile:function(oEvent){/* TODO insert action*/},
		    renameFile:function(oEvent){/* TODO insert action*/},
		    uploadFile:function(oEvent){/* TODO insert action*/},
		    saveClicked:function(oEvent){/* TODO insert saving action*/ dialog.close();},
		    cancelClicked:function(oEvent){/* TODO insert cancel action*/ dialog.close();},
		    fileUploadFailed:function(oEvent){/* TODO insert action, if necessary, if not delete this line*/},
		    beforeUploadFile:function(oEvent){/* TODO insert action, if necessary, if not delete this line*/},
		    useMultipart:true,
		    renameEnabled:true,
		    showNoData:false,
		    useEditControls:true,
		    uploadEnabled:true,
		    showAttachmentsLabelInEditMode:true,
		    editMode:false
		});
		
		
		 dialog = new sap.m.Dialog({
			customHeader :new sap.m.Bar({
				contentLeft : [new sap.m.Button({
					icon : 'sap-icon://nav-back',
					press : function(oEvent){ dialog.close(); }
				})],
				contentMiddle : [new sap.m.Label({
					text : "{i18n>DIALOGS_FORWARD_ATTACHMENTS}"
				})]
			}) ,
			contentHeight:"317px",
			content : [oFileUpload],
			
		});
		dialog.placeAt(this.getView());
		dialog.open();
	}
	
});