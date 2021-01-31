const puppeteer = require('puppeteer');

function sleep(time) {
    return new Promise((resolve, reject) => {
        setTimeout(resolve, time);
    });
}

const email = '';
const password = '';

(async () => {
    const browser = await puppeteer.launch({
        headless: false,
        slowMo: 30
    });

    const page = (await browser.pages())[0];
    await page.goto('https://cookpad.com/identity/session/new');
    await page.type('input[name="identifier"]', email);
    await page.type('input[type="password"]', password);
    await page.click('input[type="submit"]');
    console.log('Please login to cookpad.');

    await page.waitForNavigation();
    console.log('Start scraping');
    await page.goto('https://cookpad.com/honor_recipe/all');
    while (true) {
        await sleep(1000);
        const links = await page.$$eval('a.recipe-title', elements => elements.map(element => element.getAttribute('href')));
        console.log(links);
        await page.click('a.next_page');
    }
})();