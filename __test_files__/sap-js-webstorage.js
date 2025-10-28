//jQuery.sap.require("sap.ca.scfld.md.controller.ScfldMasterController");
//jQuery.sap.require("i2d.qm.qualityissue.confirm.utils.Helper");
//jQuery.sap.require("i2d.qm.qualityissue.confirm.utils.FragmentHelper");
//jQuery.sap.require("sap.ca.ui.DatePicker");
//jQuery.sap.require("sap.ca.scfld.md.app.MasterHeaderFooterHelper");
//jQuery.sap.require("i2d.qm.qualityissue.confirm.utils.StatusHelper");

sap.ca.scfld.md.controller.ScfldMasterController.extend("i2d.qm.qualityissue.confirm.view.S2", {
	/**
	 * @override
	 * 
	 * Called by the UI5 runtime to init this controller
	 * 
	 */
	onInit : function() {
		// Execute onInit for the base class
		// BaseMasterController
		sap.ca.scfld.md.controller.ScfldMasterController.prototype.onInit.call(this);
		// Settings
		this.oMasterModel = new sap.ui.model.json.JSONModel({
			selectedFilter : "All",
			selectedSorter : "CreatedOn",
			toogleSubmit : false
		});
		this.getView().setModel(this.oMasterModel, "masterModel");
		
		// Retrieve the application bundle
		this.resourceBundle = this.oApplicationFacade.getResourceBundle();
		this.oSettingsDialog = null;
		// Set list max count displayed
		this.SETTINGS_NAME = sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().settingsName;
		this.objSettings = localStorage.getObj(this.SETTINGS_NAME);
		
	 	//CSS Fix -remove leading zeros coming from oData
		if (this.objSettings.maxHits){
			this.objSettings.maxHits = parseInt(this.objSettings.maxHits, 10);
			if (isNaN(this.objSettings.maxHits)){
				jQuery.sap.log.error("Error with data, please check consistency of record. Default 30 issues displayed.");
				this.objSettings.maxHits = '30';
			}
		} //end of fix
		
		if ($.isBlank(this.objSettings)) {
			var settingsMapping = [
			                       {output:"maxHits", source:"MaxHits"},
			                       {output:"defPlant", source:"Plant"},
			                       {output:"maxFileSize", source:"MaxFileSize"}
			                       ];
			var batchResult = i2d.qm.qualityissue.confirm.utils.Helper.getCollection([{indexCollection: 3, arConversionRules : settingsMapping}], this);
			this.objSettings = batchResult && batchResult.length >0 && batchResult[0];
			// store the object into local Storage for future usage
			localStorage.setObj(this.SETTINGS_NAME, this.objSettings);
		}
		
		
		this.setMaxHitsToList();

		var bus = sap.ui.getCore().getEventBus();				   
        
		bus.subscribe(sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getAppIdentifier(), "RefreshDetail", this.handleRefresh, this);

		var oList = this.getList();
		oList.attachUpdateFinished(this.handleUpdateFinished, this);
		
		var oTemplate = oList.getItems()[0].clone();
		oList.bindItems("/QMNotificationSelectionSet", oTemplate);
		
		this.aFilterBy = [];
		this.dateFromUTC;
		this.dateToUTC;
		this.DateFromUTCValue;
		this.DateToUTCValue;
		this.previousDateFromUTC;
		this.previousDateToUTC;
		
		this.oRouter.attachRouteMatched(function(oEvent) {
			if (oEvent.getParameter("name") === "masterDetail" && oEvent.getParameter("arguments").id === "create" && this.getList().getSelectedItem()) {
				// if creating for a second time (pressing plus button on the list) then the list should be updated  
				// and the selection should be always the first line in list 
				this.handleRefresh();
			}
			}, this);
		
	},
	
	navToEmptyView : function(){
		// work around until release 1.16.6, which will contains fix for this bug
		this.showEmptyView("DETAIL_TITLE", "NO_ITEMS_AVAILABLE"); 
	}, 
	
	handleUpdateFinished : function(oControlEvent) {

		if (!jQuery.device.is.phone && (oControlEvent.getParameters().reason === "Filter"  || oControlEvent.getParameters().reason === "Change" || 
										oControlEvent.getParameters().reason === "Refresh" || oControlEvent.getParameters().reason === "Binding")) {
			
			//this._selectDetail();

			if (this.oRouter._oRouter._prevMatchedRequest === "noData/DETAIL_TITLE/NO_ITEMS_AVAILABLE") {
				this.getList().removeSelections();
			}

		}
	},

	handleRefresh : function() {
		this.setMaxHitsToList();
		this.getList().removeSelections();
		this.getList().getBinding("items")._refresh();
	},

	getHeaderFooterOptions : function() {
		return {
			sI18NMasterTitle : "QI_TITLE_MASTER_VIEW",
			buttonList : [],

			aAdditionalSettingButtons : [ {
				sId : "Settings",
				sI18nBtnTxt : "QI_MV_SETTINGS",
				sIcon : "sap-icon://wrench",
				onBtnPressed : jQuery.proxy(function(oEvent) {
					this.onSettingsPressed();
				}, this)
			} ],

			oFilterOptions : {
				onFilterPressed : $.proxy(this.onFilter, this)
			},
			oSortOptions : {
				onSortPressed : $.proxy(this.onSort, this)
			},
			onAddPress : jQuery.proxy(function(evt) {
				this.onCreate();
				jQuery.sap.log.info("add pressed");
			}, this)
		};
	},

	setMaxHitsToList : function() {
		// Evil hack until the solution is provided
		var oModel = this.getList().getModel();
		oModel.setCountSupported(false);
		oModel.setSizeLimit(this.objSettings.maxHits);
	},

	
	/**
	 * @override
	 * 
	 * @param oItem
	 * @param sFilterPattern
	 * @returns {*}
	 */
	applySearchPatternToListItem : function(oItem, sFilterPattern) {
		if (oItem.isSelectable()) {
			if (sFilterPattern.substring(0, 1) === "#") {
				var sTail = sFilterPattern.substr(1);
				var sDescr = oItem.getBindingContext().getProperty("Name").toLowerCase();
				return sDescr.indexOf(sTail) === 0;
			} else {
				return sap.ca.scfld.md.controller.ScfldMasterController.prototype.applySearchPatternToListItem.call(null, oItem, sFilterPattern);
			}
		}

	},

	/**
	 * @override
	 * 
	 * determines whether search is triggered with each change of the search
	 * field content (or only when the user explicitly starts the search).
	 * 
	 */

	isLiveSearch : function() {
		return false;
	},

	/**
	 * Called by the UI5 runtime to cleanup this controller
	 */

	onExit : function() {
		
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.destroySortDialog();
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.destroyFilterDialog();
		
		// destroy the control if needed
		if (this.oSettingsDialog) {
			this.oSettingsDialog.destroy();
			this.oSettingsDialog = null;
			this.plantlist.destroy();
			this.plantlist = null;
			this.oSimpleForm.destroy();
			this.oSimpleForm = null;
			this.settingsLists.destroy();
			this.settingsLists = null;
		}
		
		if (this.oMasterModel) {
			this.oMasterModel.destroy();
			this.oMasterModel = null;
		}
		
	},

	setInfoLabel : function(resourceId, sText, bInfoEnabled, sFilteredBy) {
		if (bInfoEnabled === null)
			bInfoEnabled = true;

		this.oMasterModel.setProperty('/toogleSubmit', bInfoEnabled);
		if (bInfoEnabled === false)
			return;

		var oLabelToolbar = this.getView().byId("labelTB");
		var toReplace = "";
		if (sText)
			toReplace = this.resourceBundle.getText(sText);
		else
			toReplace = sFilteredBy;
		var sText = this.resourceBundle.getText(resourceId, [ toReplace ]);
		oLabelToolbar.setText(sText);
		oLabelToolbar.setTooltip(sText);
	},	
	
	setFilterInfoLabel : function(arFilterBy) {

		if (!$.isArray(arFilterBy) || arFilterBy.length < 1 ){
			this.setInfoLabel('', '', false, '');
			return;
		}
		
		var infoLabelText 	= "", 
			reportOn 		= this.resourceBundle.getText("QI_REPORT_ON"),
			status 			= this.resourceBundle.getText("QI_STATUS_TEXT");
		
		$.each(arFilterBy, function(index, filterBy) {

			if (self.i2d.qm.qualityissue.confirm.utils.Helper.isValidDate(filterBy.oValue1)){
				
				if (filterBy.oValue1 === sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultFromDate && 
						filterBy.oValue2 === sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultToDate){
					return false;
				}
				
				filterBy.oValue1 = filterBy.oValue1.substring(0, 10);
				filterBy.oValue2 = filterBy.oValue2.substring(0, 10);
				infoLabelText += reportOn + ": " + filterBy.oValue1 + ", " + filterBy.oValue2;
				
			}else if(filterBy.oValue1){
				if (infoLabelText.indexOf(status) === -1){
					infoLabelText += status + ": " + i2d.qm.qualityissue.confirm.utils.StatusHelper.getStatusText(filterBy.oValue1);
				}else {
					infoLabelText += ", " + i2d.qm.qualityissue.confirm.utils.StatusHelper.getStatusText(filterBy.oValue1);
				}
			}
			
		});
		
		infoLabelText += ";";
		
		this.setInfoLabel("QI_FILTERED_BY", null, true, infoLabelText);
	},

	onCreate : function() {
		this.oRouter.navTo("fsS4");

	},

	/**
	 * @override
	 * 
	 * call in the BaseMasterController
	 *  
	 */

	// Settings dialog
	onSettingsPressed : function() {
		// max hit Form

		if (!this.oSettingsDialog) {
			var oLabel = new sap.m.Label("lbl1", {
				text : this.resourceBundle.getText("QI_NUMBER_ISSUES")
			});
			var omaxHit = new sap.m.Input("maxhit", {
				showValueHelp : false,
				value : this.objSettings.maxHits,
				type : sap.m.InputType.Number
			});

			this.oSimpleForm = new sap.ui.layout.form.SimpleForm({
				editable : true,
				labelMinWidth : 30,
				layout : sap.ui.layout.form.SimpleFormLayout.ResponsiveGridLayout,
				content : [ oLabel, omaxHit, ]
			});

			// plants List
			var arSettingfProperties = [ {
				output : "city",
				source : "City"
			}, {
				output : "postCode",
				source : "PostCode"
			}, {
				output : "name",
				source : "Name"
			}, {
				output : "plant",
				source : "Plant"
			}, ];
			var batchResult = i2d.qm.qualityissue.confirm.utils.Helper.getCollection([{indexCollection: 2, arConversionRules : arSettingfProperties, itemsPrefix: "items"}], this);
			var listData = batchResult && batchResult.length >0 && batchResult[0];
			
			var formattedListData = [];
			$.each(listData.items, function(index, listItem) {
				var obj = {};
				obj.plant = listItem.plant;
				obj.description = listItem.postCode + " " + listItem.city + " " + listItem.name;
				formattedListData.push(obj);
			});

			formattedListData = {
				items : formattedListData
			};
			var oModelList = new sap.ui.model.json.JSONModel();
			oModelList.setData(formattedListData);
			var plant = "";

			var itemTemplate = new sap.m.StandardListItem({
				title : "{plant}",
				description : "{description}",
				type : sap.m.ListType.Active
			});

			this.plantlist = new sap.m.List("Plant", {
				mode : sap.m.ListMode.SingleSelectLeft,
				// Plant is selected
				select : function(oEvent) {
					var item = oEvent.getParameter("listItem");
					plant = item.getProperty("title");

				}
			});
			this.plantlist.setModel(oModelList);
			this.plantlist.bindItems("/items", itemTemplate);
			this.plantlist.setGrowing(true);

			var listItems = this.plantlist.getItems();

			for ( var i = 0; i < listItems.length; i++) {
				if (listItems[i].getProperty("title") === this.objSettings.defPlant) {
					this.plantlist.setSelectedItem(listItems[i]);
					break;
				}
			}

			// Define custom header with Navigation Back button
			var oCustomHeader = new sap.m.Bar({
				contentLeft : [ new sap.m.Button({
					icon : "sap-icon://nav-back",
					press : jQuery.proxy(function(oEvent) {
						if (this.onValidateNumberInput(omaxHit, this)) {
							this.objSettings.maxHits = omaxHit.getValue();
							if (this.objSettings.maxHits === "1") {
								this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES_ONE"));
							}								
							else {
								this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES", this.objSettings.maxHits));
							}
															
							// Empty all the content of the dialog
							this.oSettingsDialog.removeAllContent();
							// put settings List as content
							this.oSettingsDialog.addContent(this.settingsLists);
							// Remove the custom header
							this.oSettingsDialog.setCustomHeader();
						}
					}, this)
				}) ],
				contentMiddle : [ new sap.m.Label({
					text : this.resourceBundle.getText("QI_MV_SETTINGS")
				}) ],
				contentRight : []
			});

			// settings List
			this.settingsLists = new sap.m.List("settingsLists", {
				items : [ new sap.m.StandardListItem({
					title : (this.objSettings.maxHits == "1") ? this.resourceBundle.getText("QI_MAX_ISSUES_ONE") : this.resourceBundle.getText("QI_MAX_ISSUES", this.objSettings.maxHits),
					type : sap.m.ListType.Navigation,
					press : jQuery.proxy(function() {
						// remove the settings List
						this.oSettingsDialog.removeContent(this.settingsLists);
						// set the max hit  List as content
						this.oSettingsDialog.addContent(this.oSimpleForm);
						// set custom header
						// with back button
						this.oSettingsDialog.setCustomHeader(oCustomHeader);
					}, this)
				}),

				new sap.m.StandardListItem({
					title : this.resourceBundle.getText("QI_PLANT"),
					type : sap.m.ListType.Navigation,
					press : jQuery.proxy(function() {
						// remove the
						// settings List
						this.oSettingsDialog.removeContent(this.settingsLists);
						// set plant List as
						// content
						this.oSettingsDialog.addContent(this.plantlist);
						// set custom header
						// with back button
						this.oSettingsDialog.setCustomHeader(oCustomHeader);
					}, this)
				}), ]
			});

			// Dialog
			this.oSettingsDialog = new sap.m.Dialog({
				title : this.resourceBundle.getText("QI_MV_SETTINGS"),
				stretchOnPhone : true,
				content : [ this.settingsLists ],
				beginButton : new sap.m.Button({
					text : this.resourceBundle.getText("QI_CONFIRM_BTN"),
					type : sap.m.ButtonType.Default,
					press : jQuery.proxy(function() {
						if (this.onValidateNumberInput(omaxHit, this)) {
							// Set plant value
							if (plant !== "")
								this.objSettings.defPlant = plant;

							var doRefresh = false;
							if (localStorage.getObj(this.SETTINGS_NAME).maxHits !== omaxHit.getValue()) {
								doRefresh = true;
								// Set max hits value
								this.objSettings.maxHits = omaxHit.getValue();
								if (this.objSettings.maxHits === "1")
									this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES_ONE"));
								else
									this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES", this.objSettings.maxHits));
							}

							// store the object into local Storage for
							// future usage
							localStorage.setObj(this.SETTINGS_NAME, this.objSettings);
							// Close the dialog
							this.oSettingsDialog.close();
							// refresh the list with the new settings
							if (doRefresh)
								this.handleRefresh();
						}
					}, this)
				}),
				endButton : new sap.m.Button({
					text : this.resourceBundle.getText("QI_CANCEL_BTN"),
					type : sap.m.ButtonType.Default,
					press : jQuery.proxy(function() {
						this.plantlist.removeSelections();
						plant = localStorage.getObj(this.SETTINGS_NAME).defPlant;
						for ( var i = 0; i < listItems.length; i++) {
							if (listItems[i].getProperty("title") === plant) {
								this.plantlist.setSelectedItem(listItems[i]);
								break;
							}
						}
						omaxHit.setValue(localStorage.getObj(this.SETTINGS_NAME).maxHits);
						omaxHit.setValueState(sap.ui.core.ValueState.None);
						if (this.objSettings.maxHits === "1")
							this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES_ONE"));
						else
							this.settingsLists.getItems()[0].setTitle(this.resourceBundle.getText("QI_MAX_ISSUES", localStorage.getObj(this.SETTINGS_NAME).maxHits));

						// Close the dialog
						this.oSettingsDialog.close();
					}, this)
				})
			});
		} else {
			// Empty all the content of the dialog
			this.oSettingsDialog.removeAllContent();
			// put settings List as content
			this.oSettingsDialog.addContent(this.settingsLists);
			// Remove the custom header
			this.oSettingsDialog.setCustomHeader();
		}
		// Open the dialog
		this.oSettingsDialog.open();

	},

	openMaxHits : function(event) {
		var hitList = new sap.m.List({
			items : [ new sap.m.InputListItem({
				label : "Max Hit",
				content : new sap.m.Slider({
					min : 0,
					max : 1000,
					value : 7,
					width : "200px"
				})
			}) ]
		});

		oDialog1.addContent(hitList);
	},
	
	onValidateNumberInput : function(oField, oController) {
		var value = oField.getValue();
		var regex = /^[0-9]+$/;
		if (regex.test(value) && value > 0) {
			oField.setValueState(sap.ui.core.ValueState.None);
			return true;
		} else {
			oField.setValueState(sap.ui.core.ValueState.Error);
			oField.setValueStateText(oController.resourceBundle.getText("QI_INVALID_ERROR"));
			return false;
		}
	},
	
	onFilter : function(event){
		this.aFilterBy = [];
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.openFilterDialog(this);
	},
	
	onConfirmFilterDialog : function(oEvent){

		var parameters 	= oEvent.getParameters(),
			self		= this;
		
		$.each(parameters.filterItems, function(index, filterItem){
			if(filterItem.sId === "StatusNew" || filterItem.sId === "StatusInProcess" || filterItem.sId === "StatusCompleted" || filterItem.sId === "StatusPostponed"){
				self.aFilterBy.push(filterItem.getCustomData()[0].getValue().filters);
			}
		});

		if (this.shouldCreateDateFilterObject()){
			this.aFilterBy.push(new sap.ui.model.Filter("CreatedOn", sap.ui.model.FilterOperator.BT, this.dateFromUTC, this.dateToUTC));
		}
		
		//Clear default values of dates if needed
		if (this.dateFromUTC === sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultFromDate){ 
			this.dateFromUTC = null;
		}
		
		if (this.dateToUTC === sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultToDate){ 
			this.dateToUTC = null;
		}

		this.getList().getBinding("items").filter(this.aFilterBy);
		this.setFilterInfoLabel(this.aFilterBy);
		
		this.previousDateFromUTC = this.DateFromUTCValue;
		this.previousDateToUTC = this.DateToUTCValue;
	},
	
	onCancelFilterDialog : function(oEvent){
		//Because the Date filter is constructed via custom controls, the date pickers and the filter counter are manually reset. Previous values are restored
		this.DateFromUTCValue = this.previousDateFromUTC;
		this.DateToUTCValue = this.previousDateToUTC;
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.resetFilterDialogDateFilter(this.previousDateFromUTC, this.previousDateToUTC);
	},
	
	onChangeDateFrom : function(oEvent){
		//Preserve previous value for the case of canceling changes
		this.DateFromUTCValue = oEvent.getSource().mProperties.value;
		this.dateFromUTC = i2d.qm.qualityissue.confirm.utils.Helper.convertToISODateTime(oEvent.getSource().mProperties.dateValue);
		this.setDateFilterCount();
	},
	
	onChangeDateTo : function(oEvent){
		//Preserve previous value for the case of canceling changes
		this.DateToUTCValue = oEvent.getSource().mProperties.value;
		this.dateToUTC = i2d.qm.qualityissue.confirm.utils.Helper.convertToISODateTime(oEvent.getSource().mProperties.dateValue);
		this.setDateFilterCount();
	},
	
	onResetFilterDialog : function(oEvent){
		//Because the Date filter is constructed via custom controls, the date pickers and the filter counter are manually reset
		this.DateFromUTCValue = null;
		this.DateToUTCValue = null;
		this.dateFromUTC = null;
		this.dateToUTC = null;
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.resetFilterDialogDateFilter();
	},
	
	shouldCreateDateFilterObject : function(){
		
		//If both dates are blank at the same time, no date filter object should be created 
		if ($.isBlank(this.dateFromUTC) && $.isBlank(this.dateToUTC)){
			return false;
		}
	
		//Check if any of the dates should be set to its default value 
		if ($.isBlank(this.dateFromUTC)){
			this.dateFromUTC = sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultFromDate;
			return true;
		}
		
		if ($.isBlank(this.dateToUTC)){
			this.dateToUTC = sap.ca.scfld.md.app.Application.getImpl().oConfiguration.getParams().filterDialogDefaultToDate;
		}
		
		return true;
	}, 
	
	setDateFilterCount : function(){
		
		if ($.isBlank(this.dateFromUTC) && $.isBlank(this.dateToUTC)){
			i2d.qm.qualityissue.confirm.utils.FragmentHelper.setFilterDialogDateFilterCount(0);
			return;
		}
		
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.setFilterDialogDateFilterCount(1);
	},
	
	onSort : function(oEvent){
		i2d.qm.qualityissue.confirm.utils.FragmentHelper.openSortDialog(this);
	},
	
	onConfirmSortDialog : function(oEvent){

		var parameters = oEvent.getParameters();
		
		//If no no sorted item was selected, return
		if(!parameters.sortItem){
			return;
		}
		
		var sorter = parameters.sortItem.getCustomData()[0].getValue().sorter;
		
		if(!sorter){
			return;
		}
		
		//The event parameter for descending is true/false according to the user selection of Ascending/Descending on the popup
		sorter.bDescending = parameters.sortDescending;
		this.getList().getBinding("items").sort(sorter);
	}

});