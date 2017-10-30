var title = [];
var casper = require('casper').create({
    pageSettings: {
        loadImages: false,//The script is much faster when this field is set to false
        loadPlugins: false,
        userAgent: 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36'
    }
});
 

function getTitle() {
    var title = document.querySelectorAll('your-selector');
    return Array.prototype.map.call(title, function(a) {
        return a.getAttribute('title');
    });
};
casper.start(casper.cli.get(0), function() {
     this.echo('' + this.getCurrentUrl(), 'COMMENT');
     process.exit();
 
});
 
 
casper.run();
