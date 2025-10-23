sap.ui.core.UIComponent.extend("i2d.pp.prodorder.release.appref.Component", {

    badMethod : function () {
        var incorrectStringValueXXerrorXX = "Not localized-XXerrorXX";
        var correctLocalization = "{pseudoLocalized because of the curly bracket}";
        var recurseMeCorrect = correctLocalization;
        var incorrectRecurseMeXXerrorXX = incorrectStringValueXXerrorXX;

        // The incorrect ones
        // setText(XXerrorXX + "hi");
        // setText(i2d.pp.AtRisk + XXerrorXX);
        // setText(doing_fancy_stuff_with_operators | "XXerrorXX- too bad you used a string literal here");
        // setText("not translated-XXerrorXX");
        // setText(incorrectRecurseMeXXerrorXX);
        // setText(incorrectStringValueXXerrorXX);

        obj.setText('This is a hardcoded string');
        obj.setHeaderText('This is a hardcoded string');
        obj.setPurpose('This is a hardcoded string');



        // the ok ones

        // setText(doing_fancy_stuff_with_operators | lets_assume_you_know_what_youre_doing);
        // setText(recurseMeCorrect);
        // setText("{HI}");
        // setText(correctLocalization);

        obj.setText(i18n.someString);
        obj.setText(correctLocalization);
        obj.setText(recurseMeCorrect);


        
    }
});
