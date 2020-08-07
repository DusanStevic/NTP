import random
#outFileName=""C:\Users\Dule\Desktop\dule.txt""
#outFile=open(outFileName, "w")
#outFile.write("""Hello my name is ABCD""")
#outFile.close()

# def save(self, path, parallel, iter):
#     with open(f'{path}/{"parallel" if parallel else "serial"}{iter}.mesh', 'w') as file:
#         for i in range(self.rows):
#             for j in range(self.columns):
#                 file.write(f"{self.mesh[i * self.columns + j]} ")
#             file.write("\n")
#
# def upis_test(test):
#     csvfile = open('test_previewA.csv', 'w', newline='',encoding='utf-8')
#     obj = csv.writer(csvfile)
#     for row in test:
#         obj.writerow(row)
#         csvfile.flush()
#     csvfile.close()
if __name__ == "__main__":
    outFileName=r"C:\Users\Dule\Desktop\dule.txt"
    outFile=open(outFileName, "w")

    # outFile.close()
    for i in range(10000):
        x = random.randint(0, 500)
        y = random.randint(0, 500)
        print(str(x)+' '+str(y))
        outFile.write(str(x)+' '+str(y)+'\n')
    outFile.close()