   var dummy = "lassoHelper = $(\"<div id='lasso-selection-help' style='position:absolute;pointer-events:none;background:#cccccc;'></div>"
      
   var HTMLLayout = function() {
        // this.df = Logger.dateFormatter;
    };
    HTMLLayout.prototype = {
        getStyle : function(logLevel) {
            var style;
            if(logLevel === Logger.LEVEL.ERROR) {
                style = 'color:red';
            } else if(logLevel === Logger.LEVEL.WARN) {
                style = 'color:orange';
            } else if(logLevel === Logger.LEVEL.DEBUG) {
                style = 'color:green';
            } else if(logLevel === Logger.LEVEL.TRACE) {
                style = 'color:green';
            } else if(logLevel === Logger.LEVEL.INFO) {
                style = 'color:grey';
            } else {
                style = 'color:yellow';
            }
            return style;
        },
        format : function(logTime, logLevel, logCate, logMsg) {
            return "<div style=\"" + this.getStyle(logLevel) + "\">[" + logTime + "]" + "[" + getLevelStr(logLevel) + "][" + (logCate || "main") + "]-" + logMsg + "</div>";
        }
        
        
    };
