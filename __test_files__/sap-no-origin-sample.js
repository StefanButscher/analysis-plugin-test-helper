/*
 * Copyright (C) 2009-2013 SAP AG or an SAP affiliate company. All rights reserved
 */
jQuery.sap.declare("i2d.ps.milestones.overdue.util.originused");

i2d.ps.milestones.overdue.util.Conversions = {
//		projectNameFormatter : function(oProject){
//			var oApplicationImplementation = sap.ca.scfld.md.app.Application.getImpl();
//			var oResourceBundle = oApplicationImplementation.getResourceBundle();
//			return oResourceBundle.getText("Project") + " " + oProject;
//		},
//		
		urlFormatter : function(){
			var currentUrl = window.location.origin;
			var origin = location.origin;
			return currentUrl;
		}
};