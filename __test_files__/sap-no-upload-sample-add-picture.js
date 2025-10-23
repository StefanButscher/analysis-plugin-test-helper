jQuery.sap.declare("i2d.qm.qualityissue.confirm.control.AddPicture");
jQuery.sap.require("sap.ca.ui.AddPicture");
jQuery.sap.require('sap.ca.ui.PictureItem'); 
jQuery.sap.require("i2d.qm.qualityissue.confirm.utils.Helper");


sap.ca.ui.AddPicture.extend("i2d.qm.qualityissue.confirm.control.AddPicture", {
	
	
	_readFile : function(f) {
		var maxFileSize = null;
		var oBundle = sap.ca.scfld.md.app.Application.getImpl().getResourceBundle();
		// Check if data is initial		
		if ( maxFileSize === null){
			// try to get it from localStorage first
			//var objSettings = localStorage.getObj(sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().settingsName);
			//maxFileSize = objSettings.maxFileSize;			
		}

		if ((maxFileSize === null) || (maxFileSize === 0) || f.size <= maxFileSize){
			sap.ca.ui.AddPicture.prototype._readFile.call(this, f);
		
		}else{
			// clear the selected value
			var input = $("#" + this.getId() + "-capture");
			// in IE 10 this is working, but not for IE9 -> in IE9 the upload is on the server, so this code is not executed at all		
			input.wrap('<form>').closest('form').get(0).reset();
			input.unwrap();
			// Show error message
			var message = oBundle.getText("QI_MAX_FILE_SIZE", maxFileSize);
			i2d.qm.qualityissue.confirm.utils.ErrorDialog(message);			
		}
	},
	
	

});