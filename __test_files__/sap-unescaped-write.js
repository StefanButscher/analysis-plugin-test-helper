/*!
 * SAP UI development toolkit for HTML5 (SAPUI5)
 * 
 * (c) Copyright 2009-2013 SAP AG. All rights reserved
 */

/* ----------------------------------------------------------------------------------
 * Hint: This is a derived (generated) file. Changes should be done in the underlying 
 * source files only (*.control, *.js) or they will be lost after the next generation.
 * ---------------------------------------------------------------------------------- */

// Provides control sap.m.Slider.
jQuery.sap.declare("js.unescapedWrite.incorrectWrite");
jQuery.sap.require("sap.m.library");
jQuery.sap.require("sap.ui.core.Control");

sap.ui.core.Element.extend("sap.crm.MDualSliderLabel", {
    metadata : {
        properties : {
            "key" : {
                type : "any",
                defaultValue : null
            },
            "value" : {
                type : "string",
                defaultValue : null
            }
        }
    }
});

sap.ui.core.Control.extend("sap.crm.MDualSlider", {
    metadata : {

        // ---- object ----
        publicMethods : [
            // methods
            "stepUp", "stepDown" ],

        // ---- control specific ----
        library : "sap.m",
        properties : {
            "width" : {
                type : "sap.ui.core.CSSSize",
                group : "Appearance",
                defaultValue : '100%'
            },
            "enabled" : {
                type : "boolean",
                group : "Behavior",
                defaultValue : true
            },
            "visible" : {
                type : "boolean",
                group : "Appearance",
                defaultValue : true
            },
            "name" : {
                type : "string",
                group : "Misc",
                defaultValue : null
            },
            /*
             * "min" : { type : "float", group : "Data", defaultValue : 0 },
             * "max" : { type : "float", group : "Data", defaultValue : 100 },
             */
            /*
             * "step" : { type : "float", group : "Data", defaultValue : 1 },
             */
            "value" : {
                type : "float",
                group : "Data",
                defaultValue : 0
            },
            "value2" : {
                type : "float",
                group : "Data",
                defaultValue : 1
            }
        },
        aggregations : {
            units : {
                type : "sap.crm.MDualSliderLabel",
                multiple : true,
                singularName : "unit"
            }
        },
        events : {
            "change" : {},
            "liveChange" : {}
        }
    }
});

sap.crm.MDualSlider.M_EVENTS = {
    'change' : 'change',
    'liveChange' : 'liveChange'
};

// Start of sap\m\Slider.js
jQuery.sap.require("sap.ui.core.EnabledPropagator");
sap.ui.core.EnabledPropagator.apply(sap.crm.MDualSlider.prototype, [ true ]);

/* =========================================================== */
/* begin: lifecycle methods */
/* =========================================================== */

/**
 * Required adaptations before rendering.
 *
 * @private
 */
sap.crm.MDualSlider.prototype.onBeforeRendering = function() {
    var fMin = 0;
    var units = this.getUnits();
    var fMax = this.getUnits().length - 1;
    var fStep = 1;// this.getStep();
    var bMinbiggerThanMax = false, bError = false;

    /*
     * functional dependencies:
     *
     * min -> max max -> min
     *
     * max, min -> step max, min, step -> value
     *
     */

    // if the minimum is lower than or equal to the maximum, log a warning
    if (fMin >= fMax) {
        bMinbiggerThanMax = true;
        bError = true;
        jQuery.sap.log.warning("Warning: " + "Property wrong min: " + fMin
            + " >= max: " + fMax + " on ", this);
    }

    // if the step is negative or 0, set to 1 and log a warning
    if (fStep <= 0) {
        jQuery.sap.log.warning("Warning: "
            + "The step could not be negative on ", this);
        fStep = 1;

        // update the step to 1 and suppress re-rendering
        this.setProperty("step", fStep, true);
    }

    // the step can't be bigger than slider range, log a warning
    if (fStep > (fMax - fMin) && !bMinbiggerThanMax) {
        bError = true;
        jQuery.sap.log.warning("Warning: " + "Property wrong step: " + fStep
            + " > max: " + fMax + " - " + "min: " + fMin + " on ", this);
    }

    // update the value only if there aren't errors
    if (!bError) {
        this.setValue(this.getValue());
        this.setValue2(this.getValue2());
        var width = this.getValue2() - this.getValue();

        this._fValue = this._getPercentFromValue(this.getValue()) + "%";
        this._fvalue2 = this._getPercentFromValue(this.getValue2()) + "%";
        this._fvalueText = (this._getPercentFromValue(this.getValue()) - 1)
            + "%";
        this._fvalueText2 = (this._getPercentFromValue(this.getValue2()) - 2)
            + "%";
        this._fwidth = this._getPercentFromValue(width) + "%";
    }

    // flags
    this._bDisabled = !this.getEnabled();
};

/**
 * Required adaptations after rendering.
 *
 * @private
 */
sap.crm.MDualSlider.prototype.onAfterRendering = function() {

    // slider control container jQuery reference
    this._$SliderContainer = this.$();

    // slider jQuery reference
    this._$Slider = this._$SliderContainer.children(".sapcrmMsli");

    // progress indicator
    this._$ProgressIndicator = this._$Slider.children(".sapcrmMsliProgress");

    // handle jQuery reference
    this._$Handle = this._$Slider.children(".sapcrmMsliHandle");

    // handle for the left handle
    this._$Lhandle = this._$Handle.first();

    // handle for the right handle
    this._$Rhandle = this._$Handle.last();

    // after all calculations, makes the control visible
    this._$SliderContainer.css("visibility", "");
};

/* =========================================================== */
/* end: lifecycle methods */
/* =========================================================== */

/* =========================================================== */
/* begin: event handlers */
/* =========================================================== */

/**
 * Handle the touch start event happening on the slider.
 *
 * @param {jQuery.EventObject}
 *            oEvent The event object
 * @private
 */
sap.crm.MDualSlider.prototype.ontouchstart = function(oEvent) {
    var $Target = jQuery(oEvent.target);
    var fMin = 0;
    var fValue, fmax = this.getUnits().length - 1;

    oEvent.originalEvent._sapui_handledByControl = true;

    if (oEvent.targetTouches.length > 1 || this._bDisabled) { // suppress
        // multiTouch
        // events
        return;
    }

    // update the slider measures, those values may change in orientation
    // change
    this._recalculateStyles();
    this._fDiffX = this._fSliderPaddingLeft;

    // initialization
    if ($Target.attr("id") == "left") {
        this._fStartValue = this.getValue();
        this._handle_hold = "left";
        this._$Lhandle.css("background-color", "rgba(0, 124, 192, 0.3)");
        this._$Lhandle.css("border", "0.125rem solid #007cc0");

    } else if ($Target.attr("id") == "right") {
        this._fStartValue = this.getValue2();
        this._handle_hold = "right";
        this._$Rhandle.css("background-color", "rgba(0, 124, 192, 0.3)");
        this._$Rhandle.css("border", "0.125rem solid #007cc0");
    } else if ($Target.attr("id") == "mSlider_bar") {
        this._handle_hold = "bar";
        this._$Lhandle.css("background-color", "rgba(0, 124, 192, 0.3)");
        this._$Lhandle.css("border", "0.125rem solid #007cc0");
        this._$Rhandle.css("background-color", "rgba(0, 124, 192, 0.3)");
        this._$Rhandle.css("border", "0.125rem solid #007cc0");

        this.fNewValue_start = (((oEvent.targetTouches[0].pageX - this._fDiffX - this._fSliderOffsetLeft) / this._fSliderWidth) * (this
            .getUnits().length - 1 - fMin))
            + fMin;
    }

    this.fireLiveChange({
        value : this.getValue(),
        value2 : this.getValue2()
    });
};

/**
 * Handle the touch move event on the slider.
 *
 * @param {jQuery.EventObject}
 *            oEvent The event object
 * @private
 */
sap.crm.MDualSlider.prototype.ontouchmove = function(oEvent) {
    var $Target = jQuery(oEvent.target);
    var id = $Target.attr("id");

    if (this._bDisabled) {
        return;
    }

    var fMin = 0;
    var fMax = this.getUnits().length - 1;
    var fValue;

    var fNewValue = (((oEvent.targetTouches[0].pageX - this._fDiffX - this._fSliderOffsetLeft) / this._fSliderWidth) * (this
        .getUnits().length - 1 - fMin))
        + fMin;

    if (this._handle_hold == "left") {
        fValue = this.getValue();
        // validate, update the the slider value and the UI
        this.setValue(fNewValue);
    } else if (this._handle_hold == "right") {
        fValue = this.getValue2();
        // validate, update the the slider value and the UI
        this.setValue2(fNewValue);
    } else if (this._handle_hold == "bar") {

        this.fNewValue_end = (((oEvent.targetTouches[0].pageX - this._fDiffX - this._fSliderOffsetLeft) / this._fSliderWidth) * (this
            .getUnits().length - 1 - fMin))
            + fMin;
        var current_value = this.getValue();
        var current_value2 = this.getValue2();

        var newValue = current_value
            + (this.fNewValue_end - this.fNewValue_start);

        var fStep = 1;// this.getStep();
        var fModStepVal = Math.abs(newValue % fStep);
        newValue = fModStepVal * 2 >= fStep ? newValue + fStep - fModStepVal
            : newValue - fModStepVal;

        var newValue2 = current_value2
            + (this.fNewValue_end - this.fNewValue_start);

        fModStepVal = Math.abs(newValue2 % fStep);
        newValue2 = fModStepVal * 2 >= fStep ? newValue2 + fStep - fModStepVal
            : newValue2 - fModStepVal;

        if ((current_value != fMin) || (newValue > current_value)) {
            if ((current_value2 != fMax) || (newValue2 < current_value2)) {
                this.setValue(newValue);
                this.setValue2(newValue2);
            }
        }
        if (!(current_value === newValue)) {
            this.fNewValue_start = this.fNewValue_end;
        }

    }

    this.fireLiveChange({
        value : this.getValue(),
        value2 : this.getValue2()
    });
};

/**
 * Handle the touch end event on the slider.
 *
 * @private
 */
sap.crm.MDualSlider.prototype.ontouchend = function(oEvent) {
    var $Target = jQuery(oEvent.target);
    var fValue;
    var fValue2;

    if (this._handle_hold == "left") {
        fValue = this.getValue();
        fValue2 = this.getValue2();
        this._$Lhandle.css("background-color", "");
        this._$Lhandle.css("border", "");
    } else if (this._handle_hold == "right") {
        fValue = this.getValue();
        fValue2 = this.getValue2();
        this._$Rhandle.css("background-color", "");
        this._$Rhandle.css("border", "");
    } else if (this._handle_hold == "bar") {
        this._$Lhandle.css("background-color", "");
        this._$Lhandle.css("border", "");
        this._$Rhandle.css("background-color", "");
        this._$Rhandle.css("border", "");
        fValue = this.getValue();
        fValue2 = this.getValue2();
    }

    if (this._bDisabled) {
        return;
    }

    // remove active state
    if (this._fStartValue !== fValue) { // if the value if not the same
        this.fireChange({
            value : fValue,
            value2 : fValue2
        });
    }

    // remove unused properties
    delete this._fDiffX;
    delete this._fStartValue;
    this._handle_hold = null;
};

/**
 * Handle the touch cancel event on the slider.
 *
 * @private
 */
sap.crm.MDualSlider.prototype.ontouchcancel = sap.crm.MDualSlider.prototype.ontouchend;

/* ============================================================ */
/* end: event handlers */
/* ============================================================ */

/* =========================================================== */
/* begin: internal methods */
/* =========================================================== */

/**
 * Recalculate styles.
 *
 * @private
 */
sap.crm.MDualSlider.prototype._recalculateStyles = function() {

    // slider width
    this._fSliderWidth = this._$SliderContainer.width();

    // slider padding left
    this._fSliderPaddingLeft = parseFloat(this._$SliderContainer
        .css("padding-left"));

    // slider offset left
    this._fSliderOffsetLeft = this._$SliderContainer.offset().left;

    // handle width
    this._fHandleWidth = this._$Handle.width();
};

/**
 * Calculate percentage.
 *
 * @param {float}
 *            fValue the value from the slider
 * @private
 * @returns {float} percent
 */
sap.crm.MDualSlider.prototype._getPercentFromValue = function(fValue) {
    var fMin = 0;

    return ((fValue - fMin) / (this.getUnits().length - 1 - fMin)) * 100;
};

sap.crm.MDualSlider.prototype._validateN = function(n) {
    var sTypeofN = typeof n;

    if (sTypeofN === "undefined") {
        return 1; // default n = 1
    } else if (sTypeofN !== "number") {
        jQuery.sap.log.warning('Warning: n needs to be a number', this);
        return 0;
    } else if (Math.floor(n) === n && isFinite(n)) {
        return n;
    } else {
        jQuery.sap.log
            .warning('Warning: n needs to be a finite interger', this);
        return 0;
    }
};

sap.crm.MDualSlider.prototype._setValue = function(fNewValue) {
    var fMin = 0, fMax = this.getUnits().length - 1, fStep = 1, fValue = this
        .getValue(), fModStepVal = Math.abs(fNewValue % fStep), fPerVal;

    var value = this.getValue2();

    // validate the new value before arithmetic calculations
    if (typeof fNewValue !== "number" || !isFinite(fNewValue)) {
        jQuery.sap.log.error("Error:",
            '"fNewValue" needs to be a finite number of', this);
        return this;
    }

    // round the value to the nearest step
    fNewValue = fModStepVal * 2 >= fStep ? fNewValue + fStep - fModStepVal
        : fNewValue - fModStepVal;

    // validate that the value is between maximum and minimum
    fNewValue = fNewValue > fMax ? fMax : fNewValue < fMin ? fMin : fNewValue;

    if ((this._handle_hold != "bar") && (this._handle_hold != undefined)) {
        if (fNewValue >= value) {
            fNewValue = value - 1;
            return;
        }
    }

    // Floating-point in JavaScript are IEEE 64 bit values and has some problems
    // with big decimals.
    // Round the final value to 5 digits after the decimal point.
    fNewValue = Number(fNewValue.toFixed(5));

    // update the value and suppress re-rendering
    this.setProperty("value", fNewValue, true);

    var units = this.getUnits();
    this._valueText = " ";

    this._valueText = units[fNewValue].getValue();

    // if the value is the same, suppress DOM modifications and event fire
    if (fValue === this.getValue()) {
        return this;
    }

    if (this._$SliderContainer) { // after re-rendering

        fPerVal = this._getPercentFromValue(fNewValue) + "%";

        var width = this.getValue2() - this.getValue();
        var left = this._getPercentFromValue(this.getValue()) + "%";
        this._fwidth = this._getPercentFromValue(width) + "%";
        this._$ProgressIndicator[0].style.width = this._fwidth;
        this._$ProgressIndicator[0].style.left = left;

        $("#left_text").css("left", (parseInt(left) - 1) + "%");
        $("#left_text").text(this._valueText);

        // update the handle position
        this._$Handle[0].style.left = fPerVal;
        this._$Handle[0].title = this._valueText;
    }

    return this;
};

sap.crm.MDualSlider.prototype._setValue2 = function(fNewValue) {
    var fMin = 0, fMax = this.getUnits().length - 1, fStep = 1, fValue = this
        .getValue2(), fModStepVal = Math.abs(fNewValue % fStep), fPerVal;

    var value = this.getValue();

    // validate the new value before arithmetic calculations
    if (typeof fNewValue !== "number" || !isFinite(fNewValue)) {
        jQuery.sap.log.error("Error:",
            '"fNewValue" needs to be a finite number of', this);
        return this;
    }

    // round the value to the nearest step
    fNewValue = fModStepVal * 2 >= fStep ? fNewValue + fStep - fModStepVal
        : fNewValue - fModStepVal;

    // validate that the value is between maximum and minimum
    fNewValue = fNewValue > fMax ? fMax : fNewValue < fMin ? fMin : fNewValue;

    if ((this._handle_hold != "bar") && (this._handle_hold != undefined)) {
        if (fNewValue <= value) {
            fNewValue = value + 1;
            return;
        }
    }

    // Floating-point in JavaScript are IEEE 64 bit values and has some problems
    // with big decimals.
    // Round the final value to 5 digits after the decimal point.
    fNewValue = Number(fNewValue.toFixed(5));

    // update the value and suppress re-rendering
    this.setProperty("value2", fNewValue, true);

    var units = this.getUnits();

    this._valueText2 = " ";
    /*
     * for (i = 0; i < units.length; i++) { if (parseInt(units[i].getKey()) ===
     * fNewValue) { this._valueText2 = units[i].getValue(); } }
     */

    this._valueText2 = units[fNewValue].getValue();

    // if the value is the same, suppress DOM modifications and event fire
    if (fValue === this.getValue2()) {
        return this;
    }
    if (this._$SliderContainer) { // after re-rendering

        fPerVal = this._getPercentFromValue(fNewValue) + "%";

        var width = this.getValue2() - this.getValue();
        var left = this._getPercentFromValue(this.getValue2()) + "%";
        this._fwidth = this._getPercentFromValue(width) + "%";
        this._$ProgressIndicator[0].style.width = this._fwidth;
        this._$ProgressIndicator[0].style.right = left;

        if (fNewValue == fMax) {
            $("#right_text").css("left", (parseInt(left) - 2) + "%");
        } else {
            $("#right_text").css("left", (parseInt(left) - 1) + "%");
        }

        $("#right_text").text(this._valueText2);

        // update the handle position
        this._$Handle[1].style.left = fPerVal;
        this._$Handle[1].title = this._valueText2;
    }

    return this;
};

/* =========================================================== */
/* end: internal methods */
/* =========================================================== */

/* =========================================================== */
/* begin: API method */
/* =========================================================== */

sap.crm.MDualSlider.prototype.stepUp = function(n) {
    return this.setValue(this.getValue() + (this._validateN(n) * 1));
};

sap.crm.MDualSlider.prototype.stepDown = function(n) {
    return this.setValue(this.getValue() - (this._validateN(n) * 1));
};

sap.crm.MDualSlider.prototype.setValue = function(fNewValue) {

    /*
     * The first time when setValue() method is called, other properties may be
     * outdated, because the invocation order is not always the same.
     *
     * Overwriting this prototype method with an instance method after the first
     * call, will ensure correct calculations.
     *
     */
    this.setValue = this._setValue;

    // update the value and suppress re-rendering
    return this.setProperty("value", fNewValue, true);
};

sap.crm.MDualSlider.prototype.setValue2 = function(fNewValue) {

    /*
     * The first time when setValue() method is called, other properties may be
     * outdated, because the invocation order is not always the same.
     *
     * Overwriting this prototype method with an instance method after the first
     * call, will ensure correct calculations.
     *
     */
    this.setValue2 = this._setValue2;

    // update the value and suppress re-rendering
    return this.setProperty("value2", fNewValue, true);
};

/* =========================================================== */
/* end: API method */
/* =========================================================== */

/*
 * ! SAP UI development toolkit for HTML5 (SAPUI5)
 * 
 * (c) Copyright 2009-2013 SAP AG. All rights reserved
 */


/**
 * @class Slider renderer.
 * @static
 */
sap.crm.MDualSliderRenderer = {};

/**
 * Renders the HTML for the given control, using the provided
 * {@link sap.ui.core.RenderManager}.
 *
 * @param {sap.ui.core.RenderManager}
 *            oRm the RenderManager that can be used for writing to the render
 *            output buffer
 * @param {sap.ui.core.Control}
 *            oSlider an object representation of the slider that should be
 *            rendered
 */
sap.crm.MDualSliderRenderer.render = function(oRm, oSlider) {
    var fValue = oSlider.getValue();
    var fValue2 = oSlider.getValue2();

    var bEnabled = oSlider.getEnabled(), sTooltip = oSlider
        .getTooltip_AsString();

    // avoid render when not visible
    if (!oSlider.getVisible()) {
        return;
    }

    oRm.write("<div");
    oRm.addClass("sapcrmMsliCont");

    if (!bEnabled) {
        oRm.addClass("sapcrmMSliContDisabled");
    }

    oRm.addStyle("width", oSlider.getWidth());
    oRm.addStyle("visibility", "hidden");
    oRm.writeClasses();
    oRm.writeStyles();
    oRm.writeControlData(oSlider);

    if (sTooltip) {
        oRm.writeAttributeEscaped("title", sTooltip);
    }

    oRm.write(">");

    oRm.write('<div');
    oRm.addClass("sapcrmMsli");

    if (!bEnabled) {
        oRm.addClass("sapcrmMSliDisabled");
    }

    oRm.writeClasses();
    oRm.writeStyles();
    oRm.write(">");

    oRm.write('<div id="mSlider_bar" class="sapcrmMsliProgress" style="width: '
        + oSlider._fwidth + ' ; left: ' + oSlider._fValue + ' ; right: '
        + oSlider._fvalue2 + ';"></div>');

    // start render left slider handle

    oRm.write('<span id="left"');
    oRm.addClass("sapcrmMsliHandle");
    oRm.addStyle("left", oSlider._fValue);

    oRm.writeClasses();
    oRm.writeStyles();

    oRm.writeAttribute("title", oSlider._valueText);

    if (oSlider.getEnabled()) {
        oRm.writeAttribute("tabIndex", "0");
    }

    oRm.write('><span class="sapcrmMsliHandleInner"></span></span>');


    // start render right slider handle
    oRm.write('<span id="right"');
    oRm.addClass("sapcrmMsliHandle");
    oRm.addStyle("left", oSlider._fvalue2);

    oRm.writeClasses();
    oRm.writeStyles();

    oRm.writeAttribute("title", oSlider._valueText2);

    if (oSlider.getEnabled()) {
        oRm.writeAttribute("tabIndex", "0");
    }

    oRm.write('><span class="sapcrmMsliHandleInner"></span></span>');
    oRm.write('>' + oSlider.getValue2() + '</div>');
    oRm.write("</div>");

    oRm.write("</div>");
};